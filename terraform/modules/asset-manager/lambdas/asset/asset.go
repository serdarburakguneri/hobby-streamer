package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const AssetTableName = "asset"

type Asset struct {
	Id           int           `json:"id"`
	FileName     string        `json:"fileName"`
	UploadDate   string        `json:"uploadDate"`
	Status       string        `json:"status"`
	ContentType  *string       `json:"contentType,omitempty"`
	Duration     *int          `json:"duration,omitempty"`
	Resolution   *string       `json:"resolution,omitempty"`
	Title        *string       `json:"title,omitempty"`
	Description  *string       `json:"description,omitempty"`
	Storage      *StorageLocations `json:"storage,omitempty"`
	ThumbnailURL *string       `json:"thumbnailUrl,omitempty"`
	Stream       *StreamInfo   `json:"stream,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Attributes   AttributesMap `json:"attributes,omitempty"`
}

type StorageLocations struct {
	Raw        *S3Object `json:"raw,omitempty"`
	Transcoded *S3Object `json:"transcoded,omitempty"`
	Thumbnail  *S3Object `json:"thumbnail,omitempty"`
}

type S3Object struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type StreamInfo struct {
	HLS         *string `json:"hls,omitempty"`
	DASH        *string `json:"dash,omitempty"`
	DownloadURL *string `json:"downloadUrl,omitempty"`
	CdnPrefix   *string `json:"cdnPrefix,omitempty"`
}

//for pagination
type AssetPage struct {
	Items            []Asset
	LastEvaluatedKey map[string]types.AttributeValue
}

type AttributesMap map[string]interface{}

func (m AttributesMap) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	if m == nil {
		return &types.AttributeValueMemberNULL{Value: true}, nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &types.AttributeValueMemberS{Value: string(data)}, nil
}

func (m *AttributesMap) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	strVal, ok := av.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("expected AttributeValueMemberS for AttributesMap")
	}
	return json.Unmarshal([]byte(strVal.Value), m)
}

func SaveAsset(ctx context.Context, db *dynamodb.Client, item Asset) error {
	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(AssetTableName),
		Item:      itemMap,
	})
	return err
}

func GetAssetByID(ctx context.Context, db *dynamodb.Client, id int) (*Asset, error) {
	resp, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(AssetTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, fmt.Errorf("asset not found")
	}

	var item Asset
	err = attributevalue.UnmarshalMap(resp.Item, &item)
	return &item, err
}

func ListAssets(ctx context.Context, db *dynamodb.Client, limit int32, startKey map[string]types.AttributeValue) (*AssetPage, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(AssetTableName),
		Limit:     aws.Int32(limit),
	}

	if startKey != nil {
		input.ExclusiveStartKey = startKey
	}

	out, err := db.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []Asset
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}

	return &AssetPage{
		Items:            items,
		LastEvaluatedKey: out.LastEvaluatedKey,
	}, nil
}

func UpdateAssetStatusByID(ctx context.Context, db *dynamodb.Client, id int, status string) error {
	_, err := db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(AssetTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
		},
		UpdateExpression:          aws.String("SET #s = :status"),
		ExpressionAttributeNames:  map[string]string{"#s": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":status": &types.AttributeValueMemberS{Value: status}},
	})
	return err
}

func UpdateAssetAnalysisDetailsByID(ctx context.Context, db *dynamodb.Client, id int, contentType *string, duration *int, resolution *string) error {
	expr := "SET "
	attrNames := map[string]string{}
	attrValues := map[string]types.AttributeValue{}

	if contentType != nil {
		attrNames["#ct"] = "contentType"
		attrValues[":ct"] = &types.AttributeValueMemberS{Value: *contentType}
		expr += "#ct = :ct, "
	}
	if duration != nil {
		attrNames["#d"] = "duration"
		attrValues[":d"] = &types.AttributeValueMemberN{Value: strconv.Itoa(*duration)}
		expr += "#d = :d, "
	}
	if resolution != nil {
		attrNames["#r"] = "resolution"
		attrValues[":r"] = &types.AttributeValueMemberS{Value: *resolution}
		expr += "#r = :r, "
	}

	expr = expr[:len(expr)-2]

	_, err := db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(AssetTableName),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)}},
		UpdateExpression:          aws.String(expr),
		ExpressionAttributeNames:  attrNames,
		ExpressionAttributeValues: attrValues,
	})
	return err
}