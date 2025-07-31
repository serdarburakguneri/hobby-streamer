package valueobjects

import (
	"errors"
	"strings"
)

type Tag struct {
	value string
}

func NewTag(value string) (*Tag, error) {
	if value == "" {
		return nil, errors.New("tag cannot be empty")
	}

	if len(value) > 50 {
		return nil, errors.New("tag too long")
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("tag cannot be empty after trimming")
	}

	return &Tag{value: trimmed}, nil
}

func (t Tag) Value() string {
	return t.value
}

func (t Tag) Equals(other Tag) bool {
	return t.value == other.value
}
