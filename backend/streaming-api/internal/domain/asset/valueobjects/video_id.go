package valueobjects

import (
	"errors"
	"regexp"
)

type VideoID struct {
	value string
}

func NewVideoID(value string) (*VideoID, error) {
	if value == "" {
		return nil, ErrInvalidVideoID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidVideoID
	}

	return &VideoID{value: value}, nil
}

func (v VideoID) Value() string {
	return v.value
}

func (v VideoID) Equals(other VideoID) bool {
	return v.value == other.value
}

var ErrInvalidVideoID = errors.New("invalid video ID")
