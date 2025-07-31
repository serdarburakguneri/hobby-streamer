package cache

import (
	"encoding/json"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
)

type DomainMarshaller struct{}

func NewDomainMarshaller() *DomainMarshaller {
	return &DomainMarshaller{}
}

func (m *DomainMarshaller) MarshalBucket(bucket *bucketentity.Bucket) ([]byte, error) {
	if bucket == nil {
		return nil, errors.NewValidationError("bucket cannot be nil", nil)
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "marshal_bucket",
			"bucketKey": bucket.Key().Value(),
		})
	}

	return data, nil
}

func (m *DomainMarshaller) UnmarshalBucket(data []byte) (*bucketentity.Bucket, error) {
	if len(data) == 0 {
		return nil, errors.NewValidationError("data cannot be empty", nil)
	}

	var bucket bucketentity.Bucket
	if err := json.Unmarshal(data, &bucket); err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "unmarshal_bucket",
		})
	}

	return &bucket, nil
}

func (m *DomainMarshaller) MarshalBuckets(buckets []*bucketentity.Bucket) ([]byte, error) {
	if buckets == nil {
		return nil, errors.NewValidationError("buckets cannot be nil", nil)
	}

	data, err := json.Marshal(buckets)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "marshal_buckets",
			"count":     len(buckets),
		})
	}

	return data, nil
}

func (m *DomainMarshaller) UnmarshalBuckets(data []byte) ([]*bucketentity.Bucket, error) {
	if len(data) == 0 {
		return nil, errors.NewValidationError("data cannot be empty", nil)
	}

	var buckets []*bucketentity.Bucket
	if err := json.Unmarshal(data, &buckets); err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "unmarshal_buckets",
		})
	}

	return buckets, nil
}

func (m *DomainMarshaller) MarshalAsset(asset *assetentity.Asset) ([]byte, error) {
	if asset == nil {
		return nil, errors.NewValidationError("asset cannot be nil", nil)
	}

	data, err := json.Marshal(asset)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "marshal_asset",
			"slug":      asset.Slug().Value(),
		})
	}

	return data, nil
}

func (m *DomainMarshaller) UnmarshalAsset(data []byte) (*assetentity.Asset, error) {
	if len(data) == 0 {
		return nil, errors.NewValidationError("data cannot be empty", nil)
	}

	var asset assetentity.Asset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "unmarshal_asset",
		})
	}

	return &asset, nil
}

func (m *DomainMarshaller) MarshalAssets(assets []*assetentity.Asset) ([]byte, error) {
	if assets == nil {
		return nil, errors.NewValidationError("assets cannot be nil", nil)
	}

	data, err := json.Marshal(assets)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "marshal_assets",
			"count":     len(assets),
		})
	}

	return data, nil
}

func (m *DomainMarshaller) UnmarshalAssets(data []byte) ([]*assetentity.Asset, error) {
	if len(data) == 0 {
		return nil, errors.NewValidationError("data cannot be empty", nil)
	}

	var assets []*assetentity.Asset
	if err := json.Unmarshal(data, &assets); err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "unmarshal_assets",
		})
	}

	return assets, nil
}

func (m *DomainMarshaller) GenerateBucketKey(key string) string {
	return fmt.Sprintf("bucket:%s", key)
}

func (m *DomainMarshaller) GenerateBucketsListKey(limit int, nextKey *string) string {
	if nextKey != nil {
		return fmt.Sprintf("buckets:list:%d:%s", limit, *nextKey)
	}
	return fmt.Sprintf("buckets:list:%d", limit)
}

func (m *DomainMarshaller) GenerateAssetKey(slug string) string {
	return fmt.Sprintf("asset:%s", slug)
}

func (m *DomainMarshaller) GenerateAssetsListKey() string {
	return "assets:list"
}
