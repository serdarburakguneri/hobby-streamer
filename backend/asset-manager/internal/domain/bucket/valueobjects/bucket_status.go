package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type BucketStatus struct {
	value string
}

func NewBucketStatus(value string) (*BucketStatus, error) {
	if value == "" {
		return nil, errors.New("bucket status cannot be empty")
	}

	if !IsValidBucketStatus(value) {
		return nil, errors.New("invalid bucket status")
	}

	return &BucketStatus{value: value}, nil
}

func (s BucketStatus) Value() string {
	return s.value
}

func (s BucketStatus) Equals(other BucketStatus) bool {
	return s.value == other.value
}

func IsValidBucketStatus(status string) bool {
	switch status {
	case constants.StatusPending, constants.StatusActive, constants.StatusFailed,
		constants.AssetStatusDraft, constants.AssetStatusPublished, constants.AssetStatusScheduled, constants.AssetStatusExpired:
		return true
	}
	return false
}
