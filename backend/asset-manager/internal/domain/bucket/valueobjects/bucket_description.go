package valueobjects

import (
	"errors"
	"strings"
)

type BucketDescription struct {
	value string
}

func NewBucketDescription(value string) (*BucketDescription, error) {
	if value == "" {
		return nil, nil
	}

	if len(value) > 500 {
		return nil, errors.New("bucket description too long")
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	return &BucketDescription{value: trimmed}, nil
}

func (d BucketDescription) Value() string {
	return d.value
}

func (d BucketDescription) Equals(other BucketDescription) bool {
	return d.value == other.value
}
