package bucket

type CreateBucketCommand struct {
	Name        string
	Key         string
	Description *string
	OwnerID     *string
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
