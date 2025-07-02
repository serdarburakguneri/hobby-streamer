package bucket

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

func (r *Repository) SaveBucket(ctx context.Context, b *Bucket) error {
	now := time.Now().UTC()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	item, err := attributevalue.MarshalMap(b)
	if err != nil {
		return fmt.Errorf("failed to marshal bucket: %w", err)
	}

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to save bucket: %w", err)
	}

	return nil
}

func (r *Repository) GetBucketByID(ctx context.Context, bucketID int) (*Bucket, error) {
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberN{Value: strconv.Itoa(bucketID)},
	}

	output, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("get bucket failed: %w", err)
	}
	if output.Item == nil {
		return nil, errors.New("bucket not found")
	}

	var b Bucket
	if err := attributevalue.UnmarshalMap(output.Item, &b); err != nil {
		return nil, fmt.Errorf("unmarshal bucket failed: %w", err)
	}

	return &b, nil
}

func (r *Repository) ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*BucketPage, error) {
	input := &dynamodb.ScanInput{
		TableName:         aws.String(r.TableName),
		Limit:             aws.Int32(int32(limit)),
		ExclusiveStartKey: lastKey,
	}

	result, err := r.Client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	var buckets []Bucket
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &buckets); err != nil {
		return nil, fmt.Errorf("unmarshal bucket list failed: %w", err)
	}

	return &BucketPage{
		Items:            buckets,
		LastEvaluatedKey: result.LastEvaluatedKey,
	}, nil
}

func (r *Repository) PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error {
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
		return fmt.Errorf("failed to patch bucket: %w", err)
	}

	return nil
}
