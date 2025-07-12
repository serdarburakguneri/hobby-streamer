package bucket

const (
	ErrIDShouldNotBeSet = "id should not be set by client"
	ErrIDMismatch       = "id mismatch"
	ErrAssetNotFound    = "asset not found in bucket"
	ErrAssetExists      = "asset already exists in bucket"
)

const (
	BucketTypeCollection = "collection"
	BucketTypePlaylist   = "playlist"
	BucketTypeCategory   = "category"
)

const (
	BucketStatusActive   = "active"
	BucketStatusInactive = "inactive"
	BucketStatusDraft    = "draft"
)
