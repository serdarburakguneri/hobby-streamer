package valueobjects

import (
	"errors"
	"regexp"
)

type ImageID struct {
	value string
}

func NewImageID(value string) (*ImageID, error) {
	if value == "" {
		return nil, ErrInvalidImageID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidImageID
	}

	return &ImageID{value: value}, nil
}

func (i ImageID) Value() string {
	return i.value
}

func (i ImageID) Equals(other ImageID) bool {
	return i.value == other.value
}

var ErrInvalidImageID = errors.New("invalid image ID")
