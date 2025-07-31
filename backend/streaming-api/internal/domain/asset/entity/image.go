package entity

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

type Image struct {
	id              valueobjects.ImageID
	fileName        valueobjects.FileName
	url             string
	imageType       *valueobjects.ImageType
	storageLocation *valueobjects.S3ObjectValue
	width           *int
	height          *int
	size            *int
	contentType     *string
	streamInfo      *valueobjects.StreamInfoValue
	metadata        *string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewImage(
	id valueobjects.ImageID,
	fileName valueobjects.FileName,
	url string,
	imageType *valueobjects.ImageType,
	storageLocation *valueobjects.S3ObjectValue,
	width *int,
	height *int,
	size *int,
	contentType *string,
	streamInfo *valueobjects.StreamInfoValue,
	metadata *string,
	createdAt time.Time,
	updatedAt time.Time,
) *Image {
	return &Image{
		id:              id,
		fileName:        fileName,
		url:             url,
		imageType:       imageType,
		storageLocation: storageLocation,
		width:           width,
		height:          height,
		size:            size,
		contentType:     contentType,
		streamInfo:      streamInfo,
		metadata:        metadata,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

func (i *Image) ID() valueobjects.ImageID {
	return i.id
}

func (i *Image) FileName() valueobjects.FileName {
	return i.fileName
}

func (i *Image) URL() string {
	return i.url
}

func (i *Image) Type() *valueobjects.ImageType {
	return i.imageType
}

func (i *Image) StorageLocation() *valueobjects.S3ObjectValue {
	return i.storageLocation
}

func (i *Image) Width() *int {
	return i.width
}

func (i *Image) Height() *int {
	return i.height
}

func (i *Image) Size() *int {
	return i.size
}

func (i *Image) ContentType() *string {
	return i.contentType
}

func (i *Image) StreamInfo() *valueobjects.StreamInfoValue {
	return i.streamInfo
}

func (i *Image) Metadata() *string {
	return i.metadata
}

func (i *Image) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Image) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i *Image) IsPoster() bool {
	return i.imageType != nil && i.imageType.Value() == constants.ImageTypePoster
}

func (i *Image) IsThumbnail() bool {
	return i.imageType != nil && i.imageType.Value() == constants.ImageTypeThumbnail
}

func (i *Image) IsBanner() bool {
	return i.imageType != nil && i.imageType.Value() == constants.ImageTypeBanner
}

func (i *Image) HasStreamInfo() bool {
	return i.streamInfo != nil
}

func (i *Image) GetAspectRatio() *float64 {
	if i.width == nil || i.height == nil || *i.width == 0 || *i.height == 0 {
		return nil
	}

	ratio := float64(*i.width) / float64(*i.height)
	return &ratio
}

func (i *Image) IsLandscape() bool {
	if i.width == nil || i.height == nil {
		return false
	}
	return *i.width > *i.height
}

func (i *Image) IsPortrait() bool {
	if i.width == nil || i.height == nil {
		return false
	}
	return *i.height > *i.width
}

func (i *Image) IsSquare() bool {
	if i.width == nil || i.height == nil {
		return false
	}
	return *i.width == *i.height
}

func (i *Image) IsHighResolution() bool {
	if i.width == nil || i.height == nil {
		return false
	}
	return *i.width >= 1920 && *i.height >= 1080
}
