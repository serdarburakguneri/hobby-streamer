package events

import (
	"time"
)

type BucketUpdatedEvent struct {
	bucketID      string
	updatedFields map[string]interface{}
	timestamp     time.Time
}

func NewBucketUpdatedEvent(bucketID string, updatedFields map[string]interface{}) *BucketUpdatedEvent {
	return &BucketUpdatedEvent{
		bucketID:      bucketID,
		updatedFields: updatedFields,
		timestamp:     time.Now().UTC(),
	}
}

func (e *BucketUpdatedEvent) EventType() string {
	return "bucket.updated"
}

func (e *BucketUpdatedEvent) BucketID() string {
	return e.bucketID
}

func (e *BucketUpdatedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *BucketUpdatedEvent) UpdatedFields() map[string]interface{} {
	return e.updatedFields
}
