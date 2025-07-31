package valueobjects

import (
	"errors"
)

type BucketStatus struct {
	value string
}

func NewBucketStatus(value string) (*BucketStatus, error) {
	if value == "" {
		return nil, ErrInvalidBucketStatus
	}

	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"draft":    true,
		"archived": true,
	}

	if !validStatuses[value] {
		return nil, ErrInvalidBucketStatus
	}

	return &BucketStatus{value: value}, nil
}

func (s BucketStatus) Value() string {
	return s.value
}

func (s BucketStatus) Equals(other BucketStatus) bool {
	return s.value == other.value
}

var ErrInvalidBucketStatus = errors.New("invalid bucket status")
