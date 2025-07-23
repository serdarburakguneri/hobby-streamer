package bucket

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
)

type CreateBucketCommand struct {
	Name        string
	Key         string
	Description *string
	OwnerID     *string
	Status      *string
}

type GetBucketCommand struct {
	ID string
}

type GetBucketByKeyCommand struct {
	Key string
}

type UpdateBucketCommand struct {
	ID          string
	Name        *string
	Description *string
	OwnerID     *string
	Metadata    map[string]interface{}
	Status      *string
}

type DeleteBucketCommand struct {
	ID      string
	OwnerID string
}

type ListBucketsCommand struct {
	Limit   *int
	LastKey map[string]interface{}
}

type SearchBucketsCommand struct {
	Query   string
	Limit   *int
	LastKey map[string]interface{}
}

type GetBucketsByOwnerCommand struct {
	OwnerID string
	Limit   *int
	LastKey map[string]interface{}
}

type AddAssetToBucketCommand struct {
	BucketID string
	AssetID  string
	OwnerID  string
}

type RemoveAssetFromBucketCommand struct {
	BucketID string
	AssetID  string
	OwnerID  string
}

type GetBucketAssetsCommand struct {
	BucketID string
	Limit    *int
	LastKey  map[string]interface{}
}

func (c GetBucketCommand) ToDomainBucketID() (*bucket.BucketID, error) {
	return bucket.NewBucketID(c.ID)
}

func (c UpdateBucketCommand) ToDomainBucketID() (*bucket.BucketID, error) {
	return bucket.NewBucketID(c.ID)
}

func (c DeleteBucketCommand) ToDomainBucketID() (*bucket.BucketID, error) {
	return bucket.NewBucketID(c.ID)
}

func (c AddAssetToBucketCommand) ToDomainBucketID() (*bucket.BucketID, error) {
	return bucket.NewBucketID(c.BucketID)
}

func (c RemoveAssetFromBucketCommand) ToDomainBucketID() (*bucket.BucketID, error) {
	return bucket.NewBucketID(c.BucketID)
}

func (c GetBucketAssetsCommand) ToDomainBucketID() (*bucket.BucketID, error) {
	return bucket.NewBucketID(c.BucketID)
}
