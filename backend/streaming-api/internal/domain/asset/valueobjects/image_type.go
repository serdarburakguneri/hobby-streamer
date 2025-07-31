package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type ImageType struct {
	value string
}

func NewImageType(value string) (*ImageType, error) {
	if value == "" {
		return nil, ErrInvalidImageType
	}

	validTypes := map[string]bool{
		constants.ImageTypePoster:    true,
		constants.ImageTypeThumbnail: true,
		constants.ImageTypeBanner:    true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidImageType
	}

	return &ImageType{value: value}, nil
}

func (i ImageType) Value() string {
	return i.value
}

func (i ImageType) Equals(other ImageType) bool {
	return i.value == other.value
}

var ErrInvalidImageType = errors.New("invalid image type")
