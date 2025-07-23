package bucket

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, bucket *Bucket) error
	GetByID(ctx context.Context, id BucketID) (*Bucket, error)
	GetByKey(ctx context.Context, key string) (*Bucket, error)
	Update(ctx context.Context, bucket *Bucket) error
	Delete(ctx context.Context, id BucketID) error
	List(ctx context.Context, limit *int, lastKey map[string]interface{}) (*BucketPage, error)
	Search(ctx context.Context, query string, limit *int, lastKey map[string]interface{}) (*BucketPage, error)
	GetByOwnerID(ctx context.Context, ownerID string, limit *int, lastKey map[string]interface{}) (*BucketPage, error)
	AddAsset(ctx context.Context, bucketID BucketID, assetID string) error
	RemoveAsset(ctx context.Context, bucketID BucketID, assetID string) error
	GetAssetIDs(ctx context.Context, bucketID BucketID, limit *int, lastKey map[string]interface{}) ([]string, error)
	HasAsset(ctx context.Context, bucketID BucketID, assetID string) (bool, error)
	AssetCount(ctx context.Context, bucketID BucketID) (int, error)
}
