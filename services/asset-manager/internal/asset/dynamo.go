package asset

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (a *Asset) SetTimestamps(now string) {
	if a.UploadDate == "" {
		a.UploadDate = now
	}
}

func (a *Asset) ToDynamoAttributes() (map[string]types.AttributeValue, error) {
	attrs := map[string]types.AttributeValue{
		"id":         &types.AttributeValueMemberN{Value: strconv.Itoa(a.Id)},
		"uploadDate": &types.AttributeValueMemberS{Value: a.UploadDate},
		"status":     &types.AttributeValueMemberS{Value: a.Status},
	}

	// Optional simple fields
	addString := func(key string, val *string) {
		if val != nil {
			attrs[key] = &types.AttributeValueMemberS{Value: *val}
		}
	}
	addString("title", a.Title)
	addString("description", a.Description)
	addString("thumbnailUrl", a.ThumbnailURL)

	// Optional complex fields
	addJSON := func(key string, val interface{}) {
		if val == nil {
			return
		}
		if b, err := json.Marshal(val); err == nil {
			attrs[key] = &types.AttributeValueMemberS{Value: string(b)}
		}
	}

	addJSON("tags", a.Tags)
	addJSON("attributes", a.Attributes)
	addJSON("variants", a.Variants)

	return attrs, nil
}

func FromDynamoAttributes(attrs map[string]types.AttributeValue) (*Asset, error) {
	var a Asset

	// Required fields
	if v, ok := attrs["id"].(*types.AttributeValueMemberN); ok {
		id, err := strconv.Atoi(v.Value)
		if err != nil {
			return nil, err
		}
		a.Id = id
	}
	if v, ok := attrs["uploadDate"].(*types.AttributeValueMemberS); ok {
		a.UploadDate = v.Value
	}
	if v, ok := attrs["status"].(*types.AttributeValueMemberS); ok {
		a.Status = v.Value
	}

	// Optional simple fields
	getString := func(key string) *string {
		if v, ok := attrs[key].(*types.AttributeValueMemberS); ok {
			return &v.Value
		}
		return nil
	}
	a.Title = getString("title")
	a.Description = getString("description")
	a.ThumbnailURL = getString("thumbnailUrl")

	// Optional complex fields
	parseJSON := func(key string, out interface{}) {
		if v, ok := attrs[key].(*types.AttributeValueMemberS); ok {
			_ = json.Unmarshal([]byte(v.Value), out)
		}
	}
	parseJSON("tags", &a.Tags)
	parseJSON("attributes", &a.Attributes)
	parseJSON("variants", &a.Variants)

	return &a, nil
}

func NowUTCString() string {
	return time.Now().UTC().Format(time.RFC3339)
}