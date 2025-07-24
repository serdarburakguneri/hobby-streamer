package cache

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type CacheService interface {
	GetBucket(ctx context.Context, key string) (*bucket.Bucket, error)
	GetBuckets(ctx context.Context, limit int, nextKey *string) ([]*bucket.Bucket, error)
	GetAsset(ctx context.Context, slug string) (*asset.Asset, error)
	GetAssets(ctx context.Context) ([]*asset.Asset, error)
	SetBucket(ctx context.Context, bucket *bucket.Bucket) error
	SetBuckets(ctx context.Context, buckets []*bucket.Bucket) error
	SetAsset(ctx context.Context, asset *asset.Asset) error
	SetAssets(ctx context.Context, assets []*asset.Asset) error
	InvalidateBucketCache(ctx context.Context, key string) error
	InvalidateBucketsListCache(ctx context.Context) error
	InvalidateAssetCache(ctx context.Context, slug string) error
	InvalidateAssetsListCache(ctx context.Context) error
}
