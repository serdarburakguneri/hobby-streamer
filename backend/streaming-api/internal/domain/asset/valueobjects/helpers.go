package valueobjects

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

func CreateAssetID(value string) (*ID, error) {
	return NewID(value, "asset")
}

func CreateVideoID(value string) (*ID, error) {
	return NewID(value, "video")
}

func CreateImageID(value string) (*ID, error) {
	return NewID(value, "image")
}

func CreateOwnerID(value string) (*ID, error) {
	return NewID(value, "owner")
}

func CreateBucketID(value string) (*ID, error) {
	return NewID(value, "bucket")
}

func CreateTitle(value string) (*Text, error) {
	return NewText(value, "title", 255)
}

func CreateDescription(value string) (*Text, error) {
	return NewText(value, "description", 1000)
}

func CreateSlug(value string) (*Text, error) {
	return NewText(value, "slug", 100)
}

func CreateFileName(value string) (*Text, error) {
	return NewText(value, "filename", 255)
}

func CreateBucketName(value string) (*Text, error) {
	return NewText(value, "bucket_name", 100)
}

func CreateBucketDescription(value string) (*Text, error) {
	return NewText(value, "bucket_description", 500)
}

func CreateAssetType(value string) (*Enum, error) {
	allowedTypes := make(map[string]bool)
	for assetType := range constants.AllowedAssetTypes {
		allowedTypes[assetType] = true
	}
	return NewEnum(value, "asset_type", allowedTypes)
}

func CreateVideoType(value string) (*Enum, error) {
	allowedTypes := make(map[string]bool)
	for videoType := range constants.AllowedVideoTypes {
		allowedTypes[videoType] = true
	}
	return NewEnum(value, "video_type", allowedTypes)
}

func CreateImageType(value string) (*Enum, error) {
	allowedTypes := make(map[string]bool)
	for imageType := range constants.AllowedImageTypes {
		allowedTypes[imageType] = true
	}
	return NewEnum(value, "image_type", allowedTypes)
}

func CreateVideoFormat(value string) (*Enum, error) {
	allowedFormats := make(map[string]bool)
	for format := range constants.AllowedVideoFormats {
		allowedFormats[format] = true
	}
	return NewEnum(value, "video_format", allowedFormats)
}

func CreateVideoQuality(value string) (*Enum, error) {
	allowedQualities := make(map[string]bool)
	for quality := range constants.AllowedVideoQualities {
		allowedQualities[quality] = true
	}
	return NewEnum(value, "video_quality", allowedQualities)
}

func CreateStatus(value string) (*Enum, error) {
	allowedStatuses := make(map[string]bool)
	for status := range constants.AllowedAssetStatuses {
		allowedStatuses[status] = true
	}
	for status := range constants.AllowedVideoStatuses {
		allowedStatuses[status] = true
	}
	return NewEnum(value, "status", allowedStatuses)
}

func CreateGenre(value string) (*Enum, error) {
	allowedGenres := make(map[string]bool)
	for genre := range constants.AllowedGenres {
		allowedGenres[genre] = true
	}
	return NewEnum(value, "genre", allowedGenres)
}

func CreateBucketType(value string) (*Enum, error) {
	allowedTypes := make(map[string]bool)
	for bucketType := range constants.AllowedBucketTypes {
		allowedTypes[bucketType] = true
	}
	return NewEnum(value, "bucket_type", allowedTypes)
}

func CreateBucketStatus(value string) (*Enum, error) {
	allowedStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"draft":    true,
	}
	return NewEnum(value, "bucket_status", allowedStatuses)
}
