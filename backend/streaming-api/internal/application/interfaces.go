package application

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	appasset "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/application/asset"
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
	bucketvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
)

type AssetServiceInterface interface {
	GetAsset(ctx context.Context, slug assetvalueobjects.Slug) (*assetentity.Asset, error)
	GetAssets(ctx context.Context) ([]*assetentity.Asset, error)
	GetPublicAssets(ctx context.Context) ([]*assetentity.Asset, error)
	GetAssetsByType(ctx context.Context, assetType assetvalueobjects.AssetType) ([]*assetentity.Asset, error)
	GetAssetsByGenre(ctx context.Context, genre assetvalueobjects.Genre) ([]*assetentity.Asset, error)
	GetAssetsInBucket(ctx context.Context, bucketKey bucketvalueobjects.BucketKey) ([]*assetentity.Asset, error)
	SearchAssets(ctx context.Context, query string, filters *appasset.SearchFilters) ([]*assetentity.Asset, error)
	GetStreamingInfo(ctx context.Context, slug assetvalueobjects.Slug, userID string, region string, userAge int) (*appasset.StreamingInfo, error)
	GetRecommendedAssets(ctx context.Context, slug assetvalueobjects.Slug, limit int) ([]*assetentity.Asset, error)
	GetPublishStatus(ctx context.Context, slug assetvalueobjects.Slug) (constants.PublishStatus, error)
}

type BucketServiceInterface interface {
	GetBuckets(ctx context.Context, limit int, nextKey *string) ([]*bucketentity.Bucket, error)
	GetBucket(ctx context.Context, key bucketvalueobjects.BucketKey) (*bucketentity.Bucket, error)
}
