package application

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	assetdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	bucketdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type AssetServiceInterface interface {
	GetAsset(ctx context.Context, slug assetdomain.Slug) (*assetdomain.Asset, error)
	GetAssets(ctx context.Context) ([]*assetdomain.Asset, error)
	GetPublicAssets(ctx context.Context) ([]*assetdomain.Asset, error)
	GetAssetsByType(ctx context.Context, assetType assetdomain.AssetType) ([]*assetdomain.Asset, error)
	GetAssetsByGenre(ctx context.Context, genre assetdomain.Genre) ([]*assetdomain.Asset, error)
	GetAssetsInBucket(ctx context.Context, bucketKey bucketdomain.BucketKey) ([]*assetdomain.Asset, error)
	SearchAssets(ctx context.Context, query string, filters *assetdomain.SearchFilters) ([]*assetdomain.Asset, error)
	GetStreamingInfo(ctx context.Context, slug assetdomain.Slug, userID string, region string, userAge int) (*assetdomain.StreamingInfo, error)
	GetRecommendedAssets(ctx context.Context, slug assetdomain.Slug, limit int) ([]*assetdomain.Asset, error)
	GetPublishStatus(ctx context.Context, slug assetdomain.Slug) (constants.PublishStatus, error)
}

type BucketServiceInterface interface {
	GetBuckets(ctx context.Context) ([]*bucketdomain.Bucket, error)
	GetBucket(ctx context.Context, key bucketdomain.BucketKey) (*bucketdomain.Bucket, error)
}
