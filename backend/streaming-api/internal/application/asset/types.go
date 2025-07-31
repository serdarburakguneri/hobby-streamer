package asset

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

type SearchFilters struct {
	AssetType     *valueobjects.AssetType
	Genre         *valueobjects.Genre
	OnlyPublic    bool
	OnlyPublished bool
	OnlyReady     bool
	HasVideo      bool
	HasImage      bool
}

type StreamingInfo struct {
	AssetID      string
	Title        string
	Description  string
	VideoID      string
	VideoURL     string
	ThumbnailURL string
	Duration     *float64
	Width        *int
	Height       *int
	Format       string
	StreamInfo   *valueobjects.StreamInfoValue
}
