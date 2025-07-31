package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type VideoFormat struct {
	value string
}

func NewVideoFormat(value string) (*VideoFormat, error) {
	if value == "" {
		return nil, ErrInvalidVideoFormat
	}

	if !constants.IsValidVideoFormat(value) {
		return nil, ErrInvalidVideoFormat
	}

	return &VideoFormat{value: value}, nil
}

func (v VideoFormat) Value() string {
	return v.value
}

func (v VideoFormat) Equals(other VideoFormat) bool {
	return v.value == other.value
}

var ErrInvalidVideoFormat = errors.New("invalid video format")
