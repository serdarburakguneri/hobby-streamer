package commands

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
)

type CreateAssetCommand struct {
	Slug      valueobjects.Slug
	Title     *valueobjects.Title
	AssetType *valueobjects.AssetType
	OwnerID   *valueobjects.OwnerID
}

type PatchAssetCommand struct {
	ID      valueobjects.AssetID
	Patches []JSONPatchOperation
}

type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type DeleteAssetCommand struct {
	ID valueobjects.AssetID
}

type AddVideoCommand struct {
	AssetID         valueobjects.AssetID
	Label           string
	Format          *valueobjects.VideoFormat
	StorageLocation valueobjects.S3Object
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
	AssetID valueobjects.AssetID
	VideoID string
	//TODO: Add metadata
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
