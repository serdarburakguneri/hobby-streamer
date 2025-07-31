package events

import (
	"time"
)

type BucketStatusChangedEvent struct {
	bucketID  string
	oldStatus string
	newStatus string
	timestamp time.Time
}

func NewBucketStatusChangedEvent(bucketID, oldStatus, newStatus string) *BucketStatusChangedEvent {
	return &BucketStatusChangedEvent{
		bucketID:  bucketID,
		oldStatus: oldStatus,
		newStatus: newStatus,
		timestamp: time.Now().UTC(),
	}
}

func (e *BucketStatusChangedEvent) EventType() string {
	return "bucket.status_changed"
}

func (e *BucketStatusChangedEvent) BucketID() string {
	return e.bucketID
}

func (e *BucketStatusChangedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *BucketStatusChangedEvent) OldStatus() string {
	return e.oldStatus
}

func (e *BucketStatusChangedEvent) NewStatus() string {
	return e.newStatus
}
