package asset

import (
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type Image struct {
	ImageID          string
	ImageFileName    string
	ImageURL         string
	ImageType        ImageType
	ImageStorageLoc  *S3Object
	ImageWidth       *int
	ImageHeight      *int
	ImageSize        *int64
	ImageContentType *string
	ImageStreamInfo  *StreamInfo
	ImageMetadata    map[string]string
	ImageCreatedAt   time.Time
	ImageUpdatedAt   time.Time
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
		ImageID:         generateID(),
		ImageFileName:   fileName,
		ImageURL:        url,
		ImageType:       imageType,
		ImageStorageLoc: storageLocation,
		ImageMetadata:   make(map[string]string),
		ImageCreatedAt:  now,
		ImageUpdatedAt:  now,
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
		ImageID:          id,
		ImageFileName:    fileName,
		ImageURL:         url,
		ImageType:        imageType,
		ImageStorageLoc:  storageLocation,
		ImageWidth:       width,
		ImageHeight:      height,
		ImageSize:        size,
		ImageContentType: contentType,
		ImageStreamInfo:  streamInfo,
		ImageMetadata:    metadata,
		ImageCreatedAt:   createdAt,
		ImageUpdatedAt:   updatedAt,
	}
}

func (i Image) ID() string {
	return i.ImageID
}

func (i Image) FileName() string {
	return i.ImageFileName
}

func (i Image) URL() string {
	return i.ImageURL
}

func (i Image) Type() ImageType {
	return i.ImageType
}

func (i Image) StorageLocation() *S3Object {
	return i.ImageStorageLoc
}

func (i Image) Width() *int {
	return i.ImageWidth
}

func (i Image) Height() *int {
	return i.ImageHeight
}

func (i Image) Size() *int64 {
	return i.ImageSize
}

func (i Image) ContentType() *string {
	return i.ImageContentType
}

func (i Image) StreamInfo() *StreamInfo {
	return i.ImageStreamInfo
}

func (i Image) Metadata() map[string]string {
	return i.ImageMetadata
}

func (i Image) CreatedAt() time.Time {
	return i.ImageCreatedAt
}

func (i Image) UpdatedAt() time.Time {
	return i.ImageUpdatedAt
}

func (i *Image) SetDimensions(width, height int) error {
	if width <= 0 {
		return ErrInvalidImageWidth
	}

	if height <= 0 {
		return ErrInvalidImageHeight
	}

	i.ImageWidth = &width
	i.ImageHeight = &height
	i.ImageUpdatedAt = time.Now().UTC()
	return nil
}

func (i *Image) SetSize(size int64) error {
	if size <= 0 {
		return ErrInvalidImageSize
	}

	i.ImageSize = &size
	i.ImageUpdatedAt = time.Now().UTC()
	return nil
}

func (i *Image) SetContentType(contentType string) error {
	if contentType == "" {
		return ErrInvalidImageContentType
	}

	if len(contentType) > 100 {
		return ErrInvalidImageContentType
	}

	i.ImageContentType = &contentType
	i.ImageUpdatedAt = time.Now().UTC()
	return nil
}

func (i *Image) SetStreamInfo(streamInfo *StreamInfo) {
	i.ImageStreamInfo = streamInfo
	i.ImageUpdatedAt = time.Now().UTC()
}

func (i *Image) SetMetadata(metadata map[string]string) {
	i.ImageMetadata = metadata
	i.ImageUpdatedAt = time.Now().UTC()
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

	if i.ImageMetadata == nil {
		i.ImageMetadata = make(map[string]string)
	}

	i.ImageMetadata[key] = value
	i.ImageUpdatedAt = time.Now().UTC()
	return nil
}

func (i *Image) RemoveMetadata(key string) {
	if i.ImageMetadata != nil {
		delete(i.ImageMetadata, key)
		i.ImageUpdatedAt = time.Now().UTC()
	}
}

func (i Image) Equals(other Image) bool {
	return i.ImageID == other.ImageID &&
		i.ImageFileName == other.ImageFileName &&
		i.ImageURL == other.ImageURL &&
		i.ImageType == other.ImageType
}

var (
	ErrInvalidImageFileName      = pkgerrors.NewValidationError("invalid image file name", nil)
	ErrInvalidImageURL           = pkgerrors.NewValidationError("invalid image URL", nil)
	ErrInvalidImageWidth         = pkgerrors.NewValidationError("invalid image width", nil)
	ErrInvalidImageHeight        = pkgerrors.NewValidationError("invalid image height", nil)
	ErrInvalidImageSize          = pkgerrors.NewValidationError("invalid image size", nil)
	ErrInvalidImageContentType   = pkgerrors.NewValidationError("invalid image content type", nil)
	ErrInvalidImageMetadataKey   = pkgerrors.NewValidationError("invalid image metadata key", nil)
	ErrInvalidImageMetadataValue = pkgerrors.NewValidationError("invalid image metadata value", nil)
)
