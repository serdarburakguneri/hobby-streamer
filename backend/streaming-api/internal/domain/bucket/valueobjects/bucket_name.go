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
		return nil, ErrInvalidBucketName
	}

	if len(value) > 100 {
		return nil, ErrInvalidBucketName
	}

	return &BucketName{value: strings.TrimSpace(value)}, nil
}

func (n BucketName) Value() string {
	return n.value
}

func (n BucketName) Equals(other BucketName) bool {
	return n.value == other.value
}

var ErrInvalidBucketName = errors.New("invalid bucket name")
