package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type BucketKey struct {
	value string
}

func NewBucketKey(value string) (*BucketKey, error) {
	if value == "" {
		return nil, ErrInvalidBucketKey
	}

	if len(value) < 3 || len(value) > 50 {
		return nil, ErrInvalidBucketKey
	}

	keyRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !keyRegex.MatchString(value) {
		return nil, ErrInvalidBucketKey
	}

	return &BucketKey{value: strings.ToLower(value)}, nil
}

func (k BucketKey) Value() string {
	return k.value
}

func (k BucketKey) Equals(other BucketKey) bool {
	return k.value == other.value
}

var ErrInvalidBucketKey = errors.New("invalid bucket key")
