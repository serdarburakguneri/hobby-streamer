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
		return nil, errors.New("bucket ID cannot be empty")
	}

	if len(value) > 50 {
		return nil, errors.New("bucket ID too long")
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, value)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("bucket ID contains invalid characters")
	}

	return &BucketID{value: value}, nil
}

func (id BucketID) Value() string {
	return id.value
}

func (id BucketID) Equals(other BucketID) bool {
	return id.value == other.value
}
