package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type VideoType struct {
	value string
}

func NewVideoType(value string) (*VideoType, error) {
	if value == "" {
		return nil, ErrInvalidVideoType
	}

	validTypes := map[string]bool{
		constants.VideoTypeMain:      true,
		constants.VideoTypeTrailer:   true,
		constants.VideoTypeBehind:    true,
		constants.VideoTypeInterview: true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidVideoType
	}

	return &VideoType{value: value}, nil
}

func (v VideoType) Value() string {
	return v.value
}

func (v VideoType) Equals(other VideoType) bool {
	return v.value == other.value
}

var ErrInvalidVideoType = errors.New("invalid video type")
