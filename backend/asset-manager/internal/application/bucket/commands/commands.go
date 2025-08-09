package commands

import "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"

type CreateBucketCommand struct {
	Name        string
	Key         string
	OwnerID     *valueobjects.OwnerID
	Description *valueobjects.BucketDescription
	Type        *valueobjects.BucketType
	Status      *valueobjects.BucketStatus
	Metadata    map[string]interface{}
}

type UpdateBucketCommand struct {
	ID          valueobjects.BucketID
	Name        *valueobjects.BucketName
	Description *valueobjects.BucketDescription
	OwnerID     *valueobjects.OwnerID
	Type        *valueobjects.BucketType
	Status      *valueobjects.BucketStatus
	Metadata    map[string]interface{}
}

type DeleteBucketCommand struct {
	ID valueobjects.BucketID
}

type AddAssetToBucketCommand struct {
	BucketID valueobjects.BucketID
	AssetID  string
}

type RemoveAssetFromBucketCommand struct {
	BucketID valueobjects.BucketID
	AssetID  string
}
