package cache

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/model"
)

type CacheService interface {
	GetBucket(ctx context.Context, key string) (*model.Bucket, error)
	GetBuckets(ctx context.Context) ([]model.Bucket, error)
	GetAsset(ctx context.Context, slug string) (*model.Asset, error)
	GetAssets(ctx context.Context) ([]model.Asset, error)
	SetBucket(ctx context.Context, bucket *model.Bucket) error
	SetBuckets(ctx context.Context, buckets []model.Bucket) error
	SetAsset(ctx context.Context, asset *model.Asset) error
	SetAssets(ctx context.Context, assets []model.Asset) error
	InvalidateBucketCache(ctx context.Context, key string) error
	InvalidateBucketsListCache(ctx context.Context) error
	InvalidateAssetCache(ctx context.Context, slug string) error
	InvalidateAssetsListCache(ctx context.Context) error
}
