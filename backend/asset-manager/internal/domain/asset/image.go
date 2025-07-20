package asset

import (
	"errors"
	"time"
)

type Image struct {
	id              string
	fileName        string
	url             string
	imageType       ImageType
	storageLocation *S3Object
	width           *int
	height          *int
	size            *int64
	contentType     *string
	streamInfo      *StreamInfo
	metadata        map[string]string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewImage(
	fileName string,
	url string,
	imageType ImageType,
	storageLocation *S3Object,
) (*Image, error) {
	if fileName == "" {
		return nil, ErrInvalidImageFileName
	}

	if url == "" {
		return nil, ErrInvalidImageURL
	}

	if len(fileName) > 255 {
		return nil, ErrInvalidImageFileName
	}

	if len(url) > 2048 {
		return nil, ErrInvalidImageURL
	}

	now := time.Now().UTC()
	return &Image{
		id:              generateID(),
		fileName:        fileName,
		url:             url,
		imageType:       imageType,
		storageLocation: storageLocation,
		metadata:        make(map[string]string),
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructImage(
	id string,
	fileName string,
	url string,
	imageType ImageType,
	storageLocation *S3Object,
	width *int,
	height *int,
	size *int64,
	contentType *string,
	streamInfo *StreamInfo,
	metadata map[string]string,
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

func (i Image) ID() string {
	return i.id
}

func (i Image) FileName() string {
	return i.fileName
}

func (i Image) URL() string {
	return i.url
}

func (i Image) Type() ImageType {
	return i.imageType
}

func (i Image) StorageLocation() *S3Object {
	return i.storageLocation
}

func (i Image) Width() *int {
	return i.width
}

func (i Image) Height() *int {
	return i.height
}

func (i Image) Size() *int64 {
	return i.size
}

func (i Image) ContentType() *string {
	return i.contentType
}

func (i Image) StreamInfo() *StreamInfo {
	return i.streamInfo
}

func (i Image) Metadata() map[string]string {
	return i.metadata
}

func (i Image) CreatedAt() time.Time {
	return i.createdAt
}

func (i Image) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i *Image) SetDimensions(width, height int) error {
	if width <= 0 {
		return ErrInvalidImageWidth
	}

	if height <= 0 {
		return ErrInvalidImageHeight
	}

	i.width = &width
	i.height = &height
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Image) SetSize(size int64) error {
	if size <= 0 {
		return ErrInvalidImageSize
	}

	i.size = &size
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Image) SetContentType(contentType string) error {
	if contentType == "" {
		return ErrInvalidImageContentType
	}

	if len(contentType) > 100 {
		return ErrInvalidImageContentType
	}

	i.contentType = &contentType
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Image) SetStreamInfo(streamInfo *StreamInfo) {
	i.streamInfo = streamInfo
	i.updatedAt = time.Now().UTC()
}

func (i *Image) SetMetadata(metadata map[string]string) {
	i.metadata = metadata
	i.updatedAt = time.Now().UTC()
}

func (i *Image) AddMetadata(key, value string) error {
	if key == "" {
		return ErrInvalidImageMetadataKey
	}

	if len(key) > 50 {
		return ErrInvalidImageMetadataKey
	}

	if len(value) > 500 {
		return ErrInvalidImageMetadataValue
	}

	if i.metadata == nil {
		i.metadata = make(map[string]string)
	}

	i.metadata[key] = value
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Image) RemoveMetadata(key string) {
	if i.metadata != nil {
		delete(i.metadata, key)
		i.updatedAt = time.Now().UTC()
	}
}

func (i Image) Equals(other Image) bool {
	return i.id == other.id &&
		i.fileName == other.fileName &&
		i.url == other.url &&
		i.imageType == other.imageType
}

var (
	ErrInvalidImageFileName      = errors.New("invalid image file name")
	ErrInvalidImageURL           = errors.New("invalid image URL")
	ErrInvalidImageWidth         = errors.New("invalid image width")
	ErrInvalidImageHeight        = errors.New("invalid image height")
	ErrInvalidImageSize          = errors.New("invalid image size")
	ErrInvalidImageContentType   = errors.New("invalid image content type")
	ErrInvalidImageMetadataKey   = errors.New("invalid image metadata key")
	ErrInvalidImageMetadataValue = errors.New("invalid image metadata value")
)
