package bucket

import (
	"context"

	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
)

type Repository interface {
	GetByKey(ctx context.Context, key valueobjects.BucketKey) (*bucketentity.Bucket, error)
	GetAll(ctx context.Context, limit int, nextKey *string) ([]*bucketentity.Bucket, error)
	GetByType(ctx context.Context, bucketType valueobjects.BucketType) ([]*bucketentity.Bucket, error)
	GetActive(ctx context.Context) ([]*bucketentity.Bucket, error)
	GetByAssetType(ctx context.Context, assetType assetvalueobjects.AssetType) ([]*bucketentity.Bucket, error)
	GetByGenre(ctx context.Context, genre assetvalueobjects.Genre) ([]*bucketentity.Bucket, error)
	GetAssets(ctx context.Context, bucket *bucketentity.Bucket) ([]*assetentity.Asset, error)
}
