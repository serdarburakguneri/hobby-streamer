package asset

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	a.SetTimestamps(NowUTCString())

	item, err := a.ToDynamoAttributes()
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
		"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
	}

	out, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("get asset failed: %w", err)
	}
	if out.Item == nil {
		return nil, errors.New("asset not found")
	}

	return FromDynamoAttributes(out.Item)
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

	assets := make([]Asset, 0, len(result.Items))
	for _, item := range result.Items {
		if a, err := FromDynamoAttributes(item); err == nil {
			assets = append(assets, *a)
		}
	}

	return &AssetPage{
		Items:            assets,
		LastEvaluatedKey: result.LastEvaluatedKey,
	}, nil
}