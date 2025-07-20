package application

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	assetdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	bucketdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type AssetServiceInterface interface {
	GetAsset(ctx context.Context, slug string) (*assetdomain.Asset, error)
	GetAssets(ctx context.Context) ([]*assetdomain.Asset, error)
	GetPublicAssets(ctx context.Context) ([]*assetdomain.Asset, error)
	GetAssetsByType(ctx context.Context, assetType string) ([]*assetdomain.Asset, error)
	GetAssetsByGenre(ctx context.Context, genre string) ([]*assetdomain.Asset, error)
	GetAssetsInBucket(ctx context.Context, bucketKey string) ([]*assetdomain.Asset, error)
	SearchAssets(ctx context.Context, query string, filters *assetdomain.SearchFilters) ([]*assetdomain.Asset, error)
	GetStreamingInfo(ctx context.Context, slug string, userID string, region string, userAge int) (*assetdomain.StreamingInfo, error)
	GetRecommendedAssets(ctx context.Context, slug string, limit int) ([]*assetdomain.Asset, error)
	GetPublishStatus(ctx context.Context, slug string) (constants.PublishStatus, error)
}

type BucketServiceInterface interface {
	GetBuckets(ctx context.Context) ([]*bucketdomain.Bucket, error)
	GetBucket(ctx context.Context, key string) (*bucketdomain.Bucket, error)
}
