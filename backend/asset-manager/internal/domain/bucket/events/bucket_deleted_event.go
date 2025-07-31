package events

import (
	"time"
)

type BucketDeletedEvent struct {
	bucketID  string
	timestamp time.Time
}

func NewBucketDeletedEvent(bucketID string) *BucketDeletedEvent {
	return &BucketDeletedEvent{
		bucketID:  bucketID,
		timestamp: time.Now().UTC(),
	}
}

func (e *BucketDeletedEvent) EventType() string {
	return "bucket.deleted"
}

func (e *BucketDeletedEvent) BucketID() string {
	return e.bucketID
}

func (e *BucketDeletedEvent) Timestamp() time.Time {
	return e.timestamp
}
