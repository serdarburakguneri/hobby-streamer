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
		"fileName":   &types.AttributeValueMemberS{Value: a.FileName},
		"uploadDate": &types.AttributeValueMemberS{Value: a.UploadDate},
		"status":     &types.AttributeValueMemberS{Value: a.Status},
	}

	addJSON := func(key string, val interface{}) {
		if val == nil {
			return
		}
		if b, err := json.Marshal(val); err == nil {
			attrs[key] = &types.AttributeValueMemberS{Value: string(b)}
		}
	}

	addString := func(key string, val *string) {
		if val != nil {
			attrs[key] = &types.AttributeValueMemberS{Value: *val}
		}
	}

	addInt := func(key string, val *int) {
		if val != nil {
			attrs[key] = &types.AttributeValueMemberN{Value: strconv.Itoa(*val)}
		}
	}

	addString("contentType", a.ContentType)
	addInt("duration", a.Duration)
	addString("resolution", a.Resolution)
	addString("title", a.Title)
	addString("description", a.Description)
	addString("thumbnailUrl", a.ThumbnailURL)

	addJSON("storage", a.Storage)
	addJSON("stream", a.Stream)
	addJSON("tags", a.Tags)
	addJSON("attributes", a.Attributes)

	return attrs, nil
}

func FromDynamoAttributes(attrs map[string]types.AttributeValue) (*Asset, error) {
	var a Asset

	if v, ok := attrs["id"].(*types.AttributeValueMemberN); ok {
		id, err := strconv.Atoi(v.Value)
		if err != nil {
			return nil, err
		}
		a.Id = id
	}
	if v, ok := attrs["fileName"].(*types.AttributeValueMemberS); ok {
		a.FileName = v.Value
	}
	if v, ok := attrs["uploadDate"].(*types.AttributeValueMemberS); ok {
		a.UploadDate = v.Value
	}
	if v, ok := attrs["status"].(*types.AttributeValueMemberS); ok {
		a.Status = v.Value
	}

	unmarshalString := func(key string) *string {
		if v, ok := attrs[key].(*types.AttributeValueMemberS); ok {
			return &v.Value
		}
		return nil
	}

	unmarshalInt := func(key string) *int {
		if v, ok := attrs[key].(*types.AttributeValueMemberN); ok {
			if val, err := strconv.Atoi(v.Value); err == nil {
				return &val
			}
		}
		return nil
	}

	unmarshalJSON := func(key string, target interface{}) {
		if v, ok := attrs[key].(*types.AttributeValueMemberS); ok {
			_ = json.Unmarshal([]byte(v.Value), target)
		}
	}

	a.ContentType = unmarshalString("contentType")
	a.Duration = unmarshalInt("duration")
	a.Resolution = unmarshalString("resolution")
	a.Title = unmarshalString("title")
	a.Description = unmarshalString("description")
	a.ThumbnailURL = unmarshalString("thumbnailUrl")

	unmarshalJSON("storage", &a.Storage)
	unmarshalJSON("stream", &a.Stream)
	unmarshalJSON("tags", &a.Tags)
	unmarshalJSON("attributes", &a.Attributes)

	return &a, nil
}

func NowUTCString() string {
	return time.Now().UTC().Format(time.RFC3339)
}