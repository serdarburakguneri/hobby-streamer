package commands

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
)

type CreateAssetCommand struct {
	Slug      valueobjects.Slug
	Title     *valueobjects.Title
	AssetType *valueobjects.AssetType
	OwnerID   *valueobjects.OwnerID
	ParentID  *valueobjects.AssetID
}

type DeleteAssetCommand struct {
	ID valueobjects.AssetID
}

type AddVideoCommand struct {
	AssetID         valueobjects.AssetID
	Label           string
	Format          *valueobjects.VideoFormat
	StorageLocation valueobjects.S3Object
	StreamInfo      *valueobjects.StreamInfo
	Codec           string
	VideoCodec      string
	AudioCodec      string
	FrameRate       string
	AudioChannels   int
	AudioSampleRate int
	Duration        float64
	Bitrate         int
	Width           int
	Height          int
	Size            int64
	ContentType     string
}

type UpsertVideoCommand struct {
	AssetID            valueobjects.AssetID
	Label              string
	Format             *valueobjects.VideoFormat
	StorageLocation    valueobjects.S3Object
	StreamInfo         *valueobjects.StreamInfo
	Codec              string
	VideoCodec         string
	AudioCodec         string
	FrameRate          string
	AudioChannels      int
	AudioSampleRate    int
	Duration           float64
	Bitrate            int
	Width              int
	Height             int
	Size               int64
	ContentType        string
	InitialStatus      *valueobjects.VideoStatus
	SegmentCount       int
	AvgSegmentDuration float64
	Segments           []string
}

type RemoveVideoCommand struct {
	AssetID valueobjects.AssetID
	VideoID string
}

type UpdateVideoStatusCommand struct {
	AssetID valueobjects.AssetID
	VideoID string
	Status  valueobjects.VideoStatus
}

type UpdateVideoMetadataCommand struct {
	AssetID     valueobjects.AssetID
	VideoID     string
	Width       int
	Height      int
	Duration    float64
	Bitrate     int
	Codec       string
	Size        int64
	ContentType string
}

type AddImageCommand struct {
	AssetID valueobjects.AssetID
	Image   valueobjects.Image
}

type RemoveImageCommand struct {
	AssetID valueobjects.AssetID
	ImageID string
}

type PublishAssetCommand struct {
	AssetID     valueobjects.AssetID
	PublishRule *valueobjects.PublishRule
}

type SetAssetPublishRuleCommand struct {
	AssetID     valueobjects.AssetID
	PublishRule valueobjects.PublishRule
}

type ClearAssetPublishRuleCommand struct {
	AssetID valueobjects.AssetID
}

type UpdateAssetTitleCommand struct {
	AssetID valueobjects.AssetID
	Title   valueobjects.Title
}

type UpdateAssetDescriptionCommand struct {
	AssetID     valueobjects.AssetID
	Description valueobjects.Description
}
