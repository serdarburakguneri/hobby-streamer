package cache

import (
	"context"

	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
)

type CacheService interface {
	GetBucket(ctx context.Context, key string) (*bucketentity.Bucket, error)
	GetBuckets(ctx context.Context, limit int, nextKey *string) ([]*bucketentity.Bucket, error)
	GetAsset(ctx context.Context, slug string) (*assetentity.Asset, error)
	GetAssets(ctx context.Context) ([]*assetentity.Asset, error)
	SetBucket(ctx context.Context, bucket *bucketentity.Bucket) error
	SetBuckets(ctx context.Context, buckets []*bucketentity.Bucket, limit int, nextKey *string) error
	SetAsset(ctx context.Context, asset *assetentity.Asset) error
	SetAssets(ctx context.Context, assets []*assetentity.Asset) error
	InvalidateBucketCache(ctx context.Context, key string) error
	InvalidateBucketsListCache(ctx context.Context) error
	InvalidateAssetCache(ctx context.Context, slug string) error
	InvalidateAssetsListCache(ctx context.Context) error
}
