package valueobjects

import (
	"errors"
)

type BucketType struct {
	value string
}

func NewBucketType(value string) (*BucketType, error) {
	if value == "" {
		return nil, ErrInvalidBucketType
	}

	validTypes := map[string]bool{
		"collection": true,
		"playlist":   true,
		"category":   true,
		"featured":   true,
		"trending":   true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidBucketType
	}

	return &BucketType{value: value}, nil
}

func (t BucketType) Value() string {
	return t.value
}

func (t BucketType) Equals(other BucketType) bool {
	return t.value == other.value
}

var ErrInvalidBucketType = errors.New("invalid bucket type")
