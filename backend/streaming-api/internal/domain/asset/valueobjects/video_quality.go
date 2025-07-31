package valueobjects

import (
	"errors"
)

type VideoQuality struct {
	value string
}

func NewVideoQuality(value string) (*VideoQuality, error) {
	validQualities := map[string]bool{
		"low": true, "medium": true, "high": true, "fourk": true,
	}
	if !validQualities[value] {
		return nil, ErrInvalidVideoQuality
	}
	return &VideoQuality{value: value}, nil
}

func (q VideoQuality) Value() string                  { return q.value }
func (q VideoQuality) Equals(other VideoQuality) bool { return q.value == other.value }

var ErrInvalidVideoQuality = errors.New("invalid video quality")
