package valueobjects

import (
	"errors"
	"strings"
)

type VideoLabel struct {
	value string
}

func NewVideoLabel(value string) (*VideoLabel, error) {
	if value == "" {
		return nil, errors.New("video label cannot be empty")
	}

	if len(value) > 100 {
		return nil, errors.New("video label too long")
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("video label cannot be empty after trimming")
	}

	return &VideoLabel{value: trimmed}, nil
}

func (vl VideoLabel) Value() string {
	return vl.value
}

func (vl VideoLabel) Equals(other VideoLabel) bool {
	return vl.value == other.value
}
