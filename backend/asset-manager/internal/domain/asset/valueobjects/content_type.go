package valueobjects

import (
	"errors"
	"strings"
)

type ContentType struct {
	value string
}

func NewContentType(value string) (*ContentType, error) {
	if value == "" {
		return nil, errors.New("content type cannot be empty")
	}

	if len(value) > 100 {
		return nil, errors.New("content type too long")
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("content type cannot be empty after trimming")
	}

	if !strings.Contains(trimmed, "/") {
		return nil, errors.New("invalid content type format")
	}

	return &ContentType{value: trimmed}, nil
}

func (ct ContentType) Value() string {
	return ct.value
}

func (ct ContentType) Equals(other ContentType) bool {
	return ct.value == other.value
}

func (ct ContentType) IsVideo() bool {
	return strings.HasPrefix(ct.value, "video/")
}

func (ct ContentType) IsAudio() bool {
	return strings.HasPrefix(ct.value, "audio/")
}

func (ct ContentType) IsImage() bool {
	return strings.HasPrefix(ct.value, "image/")
}
