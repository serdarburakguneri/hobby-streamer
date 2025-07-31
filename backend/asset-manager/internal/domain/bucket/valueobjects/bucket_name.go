package valueobjects

import (
	"errors"
	"strings"
)

type BucketName struct {
	value string
}

func NewBucketName(value string) (*BucketName, error) {
	if value == "" {
		return nil, errors.New("bucket name cannot be empty")
	}

	if len(value) > 100 {
		return nil, errors.New("bucket name too long")
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New("bucket name cannot be empty after trimming")
	}

	return &BucketName{value: trimmed}, nil
}

func (n BucketName) Value() string {
	return n.value
}

func (n BucketName) Equals(other BucketName) bool {
	return n.value == other.value
}
