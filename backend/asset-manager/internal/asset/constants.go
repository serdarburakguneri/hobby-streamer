package asset

import "github.com/serdarburakguneri/hobby-streamer/pkg/constants"

const (
	// Asset-specific errors
	ErrIDShouldNotBeSet = "id should not be set by client"
	ErrIDMismatch       = "id mismatch"
	ErrImageExists      = "image with same filename already exists"

	// Asset types
	AssetTypeVideo = "video"
	AssetTypeImage = "image"
	AssetTypeAudio = "audio"

	// Asset categories
	AssetCategoryMovie       = "movie"
	AssetCategoryTVShow      = "tv_show"
	AssetCategoryDocumentary = "documentary"
	AssetCategoryMusic       = "music"
)

// Use shared constants
var (
	AssetStatusPending = constants.StatusPending
	AssetStatusActive  = constants.StatusActive
	AssetStatusFailed  = constants.StatusFailed
)
