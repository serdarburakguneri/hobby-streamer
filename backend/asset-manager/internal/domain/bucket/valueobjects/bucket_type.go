package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type BucketType struct {
	value string
}

func NewBucketType(value string) (*BucketType, error) {
	if value == "" {
		return nil, errors.New("bucket type cannot be empty")
	}

	if !IsValidBucketType(value) {
		return nil, errors.New("invalid bucket type")
	}

	return &BucketType{value: value}, nil
}

func (t BucketType) Value() string {
	return t.value
}

func (t BucketType) Equals(other BucketType) bool {
	return t.value == other.value
}

func IsValidBucketType(t string) bool {
	_, ok := constants.AllowedBucketTypes[t]
	return ok
}
