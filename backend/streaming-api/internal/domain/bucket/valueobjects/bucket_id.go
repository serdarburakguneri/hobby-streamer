package valueobjects

import (
	"errors"
	"regexp"
)

type BucketID struct {
	value string
}

func NewBucketID(value string) (*BucketID, error) {
	if value == "" {
		return nil, ErrInvalidBucketID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidBucketID
	}

	return &BucketID{value: value}, nil
}

func (id BucketID) Value() string {
	return id.value
}

func (id BucketID) Equals(other BucketID) bool {
	return id.value == other.value
}

var ErrInvalidBucketID = errors.New("invalid bucket ID")
