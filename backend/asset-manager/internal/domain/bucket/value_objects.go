package bucket

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

var ErrInvalidBucketType = pkgerrors.NewValidationError("invalid bucket type", nil)

type BucketType struct {
	value string
}

func NewBucketType(value string) (*BucketType, error) {
	if value == "" {
		return nil, ErrInvalidBucketType
	}
	if !IsValidBucketType(value) {
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

func IsValidBucketType(t string) bool {
	_, ok := constants.AllowedBucketTypes[t]
	return ok
}
