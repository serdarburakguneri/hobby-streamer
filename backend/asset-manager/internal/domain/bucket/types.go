package bucket

import (
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

var (
	ErrSlugAlreadyExists = pkgerrors.NewValidationError("bucket slug already exists", nil)
	ErrKeyAlreadyExists  = pkgerrors.NewValidationError("bucket key already exists", nil)
	ErrBucketNotFound    = pkgerrors.NewNotFoundError("bucket not found", nil)
)

type BucketPage struct {
	Items   []*Bucket              `json:"items"`
	LastKey map[string]interface{} `json:"lastKey,omitempty"`
	HasMore bool                   `json:"hasMore"`
	Total   int                    `json:"total"`
}

type BucketID struct {
	value string
}

func NewBucketID(value string) (*BucketID, error) {
	if value == "" {
		return nil, ErrBucketNotFound
	}
	// Add validation as needed (e.g., regex)
	return &BucketID{value: value}, nil
}

func (id BucketID) Value() string {
	return id.value
}

func (id BucketID) Equals(other BucketID) bool {
	return id.value == other.value
}
