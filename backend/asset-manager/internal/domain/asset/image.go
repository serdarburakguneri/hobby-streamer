package asset

import (
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type Image struct {
	id              ImageID
	fileName        ImageFileName
	url             ImageURL
	type_           ImageType
	storageLocation *S3Object
	width           *int
	height          *int
	size            *int64
	contentType     ImageContentType
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
	imageID, err := NewImageID(generateID())
	if err != nil {
		return nil, err
	}
	imageFileName, err := NewImageFileName(fileName)
	if err != nil {
		return nil, err
	}
	imageURL, err := NewImageURL(url)
	if err != nil {
		return nil, err
	}
	imageContentType, err := NewImageContentType("")
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &Image{
		id:              *imageID,
		fileName:        *imageFileName,
		url:             *imageURL,
		type_:           imageType,
		storageLocation: storageLocation,
		contentType:     *imageContentType,
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
	contentType string,
	streamInfo *StreamInfo,
	metadata map[string]string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Image, error) {
	imageID, err := NewImageID(id)
	if err != nil {
		return nil, err
	}
	imageFileName, err := NewImageFileName(fileName)
	if err != nil {
		return nil, err
	}
	imageURL, err := NewImageURL(url)
	if err != nil {
		return nil, err
	}
	imageContentType, err := NewImageContentType(contentType)
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &Image{
		id:              *imageID,
		fileName:        *imageFileName,
		url:             *imageURL,
		type_:           imageType,
		storageLocation: storageLocation,
		width:           width,
		height:          height,
		size:            size,
		contentType:     *imageContentType,
		streamInfo:      streamInfo,
		metadata:        metadata,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}, nil
}

func (i *Image) ID() ImageID {
	return i.id
}

func (i *Image) FileName() ImageFileName {
	return i.fileName
}

func (i *Image) URL() ImageURL {
	return i.url
}

func (i *Image) Type() ImageType {
	return i.type_
}

func (i *Image) StorageLocation() *S3Object {
	return i.storageLocation
}

func (i *Image) Width() *int {
	return i.width
}

func (i *Image) Height() *int {
	return i.height
}

func (i *Image) Size() *int64 {
	return i.size
}

func (i *Image) ContentType() ImageContentType {
	return i.contentType
}

func (i *Image) StreamInfo() *StreamInfo {
	return i.streamInfo
}

func (i *Image) Metadata() map[string]string {
	return i.metadata
}

func (i *Image) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Image) UpdatedAt() time.Time {
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
	newContentType, err := NewImageContentType(contentType)
	if err != nil {
		return err
	}
	i.contentType = *newContentType
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

func (i *Image) Equals(other Image) bool {
	return i.id.Equals(other.id) &&
		i.fileName.Equals(other.fileName) &&
		i.url.Equals(other.url) &&
		i.type_ == other.type_ // only compare core fields for equality
}

var (
	ErrInvalidImageWidth         = pkgerrors.NewValidationError("invalid image width", nil)
	ErrInvalidImageHeight        = pkgerrors.NewValidationError("invalid image height", nil)
	ErrInvalidImageSize          = pkgerrors.NewValidationError("invalid image size", nil)
	ErrInvalidImageContentType   = pkgerrors.NewValidationError("invalid image content type", nil)
	ErrInvalidImageMetadataKey   = pkgerrors.NewValidationError("invalid image metadata key", nil)
	ErrInvalidImageMetadataValue = pkgerrors.NewValidationError("invalid image metadata value", nil)
)
