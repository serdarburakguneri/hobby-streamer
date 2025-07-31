package bucket

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
)

type Saver interface {
	Save(ctx context.Context, bucket *entity.Bucket) error
	Update(ctx context.Context, bucket *entity.Bucket) error
	Delete(ctx context.Context, id valueobjects.BucketID) error
}

type Finder interface {
	FindByID(ctx context.Context, id valueobjects.BucketID) (*entity.Bucket, error)
	FindByKey(ctx context.Context, key valueobjects.BucketKey) (*entity.Bucket, error)
	Exists(ctx context.Context, id valueobjects.BucketID) (bool, error)
	ExistsByKey(ctx context.Context, key valueobjects.BucketKey) (bool, error)
}

type Pager interface {
	List(ctx context.Context, limit *int, offset *int) ([]*entity.Bucket, error)
	Search(ctx context.Context, query string, limit *int, offset *int) ([]*entity.Bucket, error)
	FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID, limit *int, offset *int) ([]*entity.Bucket, error)
	FindByType(ctx context.Context, bucketType valueobjects.BucketType, limit *int, offset *int) ([]*entity.Bucket, error)
	FindByStatus(ctx context.Context, status valueobjects.BucketStatus, limit *int, offset *int) ([]*entity.Bucket, error)
}

type Relation interface {
	AddAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) error
	RemoveAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) error
	GetAssetIDs(ctx context.Context, bucketID valueobjects.BucketID, limit *int, lastKey map[string]interface{}) ([]string, error)
	HasAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) (bool, error)
	AssetCount(ctx context.Context, bucketID valueobjects.BucketID) (int, error)
}

type Repository interface {
	Saver
	Finder
	Pager
	Relation
}
