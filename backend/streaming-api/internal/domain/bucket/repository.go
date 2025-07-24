package bucket

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
)

type Repository interface {
	GetByID(ctx context.Context, id BucketID) (*Bucket, error)
	GetByKey(ctx context.Context, key BucketKey) (*Bucket, error)
	GetAll(ctx context.Context, limit int, nextKey *string) ([]*Bucket, error)
	GetByType(ctx context.Context, bucketType BucketType) ([]*Bucket, error)
	GetByAssetType(ctx context.Context, assetType asset.AssetType) ([]*Bucket, error)
	GetByGenre(ctx context.Context, genre asset.Genre) ([]*Bucket, error)
	GetActive(ctx context.Context) ([]*Bucket, error)
	GetAssets(ctx context.Context, bucket *Bucket) ([]*asset.Asset, error)
	Search(ctx context.Context, query string, filters *BucketSearchFilters) ([]*Bucket, error)
	GetStats(ctx context.Context, bucket *Bucket) (*BucketStats, error)
	GetRecommendations(ctx context.Context, bucket *Bucket, limit int) ([]*asset.Asset, error)
}
