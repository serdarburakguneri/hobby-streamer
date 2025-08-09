package valueobjects

import (
	"errors"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type BucketType struct {
	value string
}

func NewBucketType(value string) (*BucketType, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return nil, ErrInvalidBucketType
	}
	if !constants.IsValidBucketType(v) {
		return nil, ErrInvalidBucketType
	}
	return &BucketType{value: v}, nil
}

func (t BucketType) Value() string {
	return t.value
}

func (t BucketType) Equals(other BucketType) bool {
	return t.value == other.value
}

var ErrInvalidBucketType = errors.New("invalid bucket type")
