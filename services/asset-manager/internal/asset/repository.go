package asset

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository struct {
	TableName string
	Client    *dynamodb.Client
}

func NewRepository(tableName string, client *dynamodb.Client) *Repository {
	return &Repository{
		TableName: tableName,
		Client:    client,
	}
}

func (r *Repository) SaveAsset(ctx context.Context, a *Asset) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if a.CreatedAt == "" {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	item, err := attributevalue.MarshalMap(a)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to save asset: %w", err)
	}

	return nil
}

func (r *Repository) GetAssetByID(ctx context.Context, id int) (*Asset, error) {
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
	}

	output, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("get asset failed: %w", err)
	}
	if output.Item == nil {
		return nil, errors.New("asset not found")
	}

	var a Asset
	if err := attributevalue.UnmarshalMap(output.Item, &a); err != nil {
		return nil, fmt.Errorf("unmarshal asset failed: %w", err)
	}

	return &a, nil
}

func (r *Repository) ListAssets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error) {
	input := &dynamodb.ScanInput{
		TableName:         aws.String(r.TableName),
		Limit:             aws.Int32(int32(limit)),
		ExclusiveStartKey: lastKey,
	}

	result, err := r.Client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	var assets []Asset
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &assets); err != nil {
		return nil, fmt.Errorf("unmarshal asset list failed: %w", err)
	}

	return &AssetPage{
		Items:            assets,
		LastEvaluatedKey: result.LastEvaluatedKey,
	}, nil
}

func (r *Repository) PatchAsset(ctx context.Context, id int, patch map[string]interface{}) error {
	if len(patch) == 0 {
		return nil
	}

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
	}

	updateExpr := "SET "
	exprAttrValues := make(map[string]types.AttributeValue)
	exprAttrNames := make(map[string]string)

	i := 0
	for k, v := range patch {
		placeholder := fmt.Sprintf("#k%d", i)
		valueKey := fmt.Sprintf(":v%d", i)

		exprAttrNames[placeholder] = k
		av, err := attributevalue.Marshal(v)
		if err != nil {
			continue // skip unsupported types
		}
		exprAttrValues[valueKey] = av

		updateExpr += fmt.Sprintf("%s = %s, ", placeholder, valueKey)
		i++
	}

	updateExpr = updateExpr[:len(updateExpr)-2] // remove trailing comma

	_, err := r.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.TableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	})
	if err != nil {
		return fmt.Errorf("failed to patch asset: %w", err)
	}

	return nil
}
