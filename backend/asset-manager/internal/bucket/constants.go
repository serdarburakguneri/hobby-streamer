package bucket

import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"

const (
	// Bucket-specific errors
	ErrIDShouldNotBeSet = "id should not be set by client"
	ErrIDMismatch       = "id mismatch"
	ErrAssetNotFound    = "asset not found in bucket"
	ErrAssetExists      = "asset already exists in bucket"

	// Bucket types
	BucketTypeCollection = "collection"
	BucketTypePlaylist   = "playlist"
	BucketTypeCategory   = "category"
)

// Use shared constants
var (
	BucketStatusPending = constants.StatusPending
	BucketStatusActive  = constants.StatusActive
	BucketStatusFailed  = constants.StatusFailed
)
