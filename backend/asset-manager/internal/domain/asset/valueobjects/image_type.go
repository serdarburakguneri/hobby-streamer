package valueobjects

import (
	"errors"
)

type ImageType string

const (
	ImageTypePoster     ImageType = "poster"
	ImageTypeBackdrop   ImageType = "backdrop"
	ImageTypeThumbnail  ImageType = "thumbnail"
	ImageTypeScreenshot ImageType = "screenshot"
	ImageTypeLogo       ImageType = "logo"
)

func NewImageType(value string) (*ImageType, error) {
	if value == "" {
		return nil, errors.New("image type cannot be empty")
	}

	validTypes := []ImageType{
		ImageTypePoster,
		ImageTypeBackdrop,
		ImageTypeThumbnail,
		ImageTypeScreenshot,
		ImageTypeLogo,
	}

	for _, imgType := range validTypes {
		if ImageType(value) == imgType {
			it := ImageType(value)
			return &it, nil
		}
	}

	return nil, errors.New("invalid image type")
}

func (it ImageType) Value() string {
	return string(it)
}

func (it ImageType) Equals(other ImageType) bool {
	return it == other
}
