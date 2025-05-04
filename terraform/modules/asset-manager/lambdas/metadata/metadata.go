package metadata

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type Metadata struct {
	Id             int                    `json:"id"`
	FileName       string                 `json:"fileName"`
	S3Bucket       string                 `json:"s3Bucket"`
	S3Key          string                 `json:"s3Key"`
	UploadDate     string                 `json:"uploadDate"`
	Status         string                 `json:"status"`
	ContentType    *string                `json:"contentType,omitempty"`
	Duration       *int                   `json:"duration,omitempty"`
	Resolution     *string                `json:"resolution,omitempty"`
	StreamURL      *string                `json:"streamUrl,omitempty"`
	ThumbnailURL   *string                `json:"thumbnailUrl,omitempty"`
	Title          *string                `json:"title,omitempty"`
	Description    *string                `json:"description,omitempty"`
}

func SaveMetadata(ctx context.Context, ddb *dynamodb.Client, tableName string, meta Metadata) error {
	item, err := attributevalue.MarshalMap(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func GetMetadata(ctx context.Context, ddb *dynamodb.Client, tableName string, id int) (*Metadata, error) {
	out, err := ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if out.Item == nil {
		return nil, fmt.Errorf("id %d not found", id)
	}

	var meta Metadata
	if err := attributevalue.UnmarshalMap(out.Item, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &meta, nil
}