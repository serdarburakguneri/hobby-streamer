package valueobjects

import (
	"errors"
	"regexp"
)

type BucketKey struct {
	value string
}

func NewBucketKey(value string) (*BucketKey, error) {
	if value == "" {
		return nil, errors.New("bucket key cannot be empty")
	}

	if len(value) < 3 || len(value) > 50 {
		return nil, errors.New("bucket key must be between 3 and 50 characters")
	}

	matched, err := regexp.MatchString(`^[a-z0-9-]+$`, value)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("bucket key contains invalid characters")
	}

	return &BucketKey{value: value}, nil
}

func (k BucketKey) Value() string {
	return k.value
}

func (k BucketKey) Equals(other BucketKey) bool {
	return k.value == other.value
}
