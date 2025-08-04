package valueobjects

import (
	"time"
)

type Image struct {
	id              ID
	fileName        ValidatedString
	url             string
	imageType       ImageType
	storageLocation *S3Object
	width           *int
	height          *int
	size            *int64
	contentType     ContentType
	streamInfo      *StreamInfo
	metadata        map[string]string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewImage(fileName, url string, imageType ImageType, contentType string) (*Image, error) {
	id, err := GenerateImageID()
	if err != nil {
		return nil, err
	}

	fileNameVO, err := NewValidatedString(fileName, 255, "fileName")
	if err != nil {
		return nil, err
	}

	contentTypeVO, err := NewContentType(contentType)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Image{
		id:          *id,
		fileName:    *fileNameVO,
		url:         url,
		imageType:   imageType,
		contentType: *contentTypeVO,
		createdAt:   now,
		updatedAt:   now,
		metadata:    make(map[string]string),
	}, nil
}

func NewImageWithDetails(
	id ID,
	fileName string,
	url string,
	imageType ImageType,
	storageLocation *S3Object,
	width, height *int,
	size *int64,
	contentType string,
	streamInfo *StreamInfo,
	metadata map[string]string,
	createdAt, updatedAt time.Time,
) (*Image, error) {
	fileNameVO, err := NewValidatedString(fileName, 255, "fileName")
	if err != nil {
		return nil, err
	}

	contentTypeVO, err := NewContentType(contentType)
	if err != nil {
		return nil, err
	}

	return &Image{
		id:              id,
		fileName:        *fileNameVO,
		url:             url,
		imageType:       imageType,
		storageLocation: storageLocation,
		width:           width,
		height:          height,
		size:            size,
		contentType:     *contentTypeVO,
		streamInfo:      streamInfo,
		metadata:        metadata,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}, nil
}

func (img *Image) ID() ID {
	return img.id
}

func (img *Image) FileName() ValidatedString {
	return img.fileName
}

func (img *Image) URL() string {
	return img.url
}

func (img *Image) Type() ImageType {
	return img.imageType
}

func (img *Image) StorageLocation() *S3Object {
	return img.storageLocation
}

func (img *Image) Width() *int {
	return img.width
}

func (img *Image) Height() *int {
	return img.height
}

func (img *Image) Size() *int64 {
	return img.size
}

func (img *Image) ContentType() ContentType {
	return img.contentType
}

func (img *Image) StreamInfo() *StreamInfo {
	return img.streamInfo
}

func (img *Image) Metadata() map[string]string {
	return img.metadata
}

func (img *Image) CreatedAt() time.Time {
	return img.createdAt
}

func (img *Image) UpdatedAt() time.Time {
	return img.updatedAt
}

func (img *Image) Equals(other Image) bool {
	return img.id.Equals(other.id)
}

func (img *Image) SetStreamInfo(streamInfo *StreamInfo) {
	img.streamInfo = streamInfo
	img.updatedAt = time.Now().UTC()
}
