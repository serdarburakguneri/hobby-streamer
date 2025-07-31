package valueobjects

import (
	"errors"
	"strings"
)

type BucketDescription struct {
	value string
}

func NewBucketDescription(value string) (*BucketDescription, error) {
	if len(value) > 500 {
		return nil, ErrInvalidBucketDescription
	}

	return &BucketDescription{value: strings.TrimSpace(value)}, nil
}

func (d BucketDescription) Value() string {
	return d.value
}

func (d BucketDescription) Equals(other BucketDescription) bool {
	return d.value == other.value
}

var ErrInvalidBucketDescription = errors.New("invalid bucket description")
