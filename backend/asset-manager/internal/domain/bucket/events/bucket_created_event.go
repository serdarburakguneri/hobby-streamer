package events

import (
	"time"
)

type BucketCreatedEvent struct {
	bucketID  string
	name      string
	key       string
	ownerID   string
	timestamp time.Time
}

func NewBucketCreatedEvent(bucketID, name, key, ownerID string) *BucketCreatedEvent {
	return &BucketCreatedEvent{
		bucketID:  bucketID,
		name:      name,
		key:       key,
		ownerID:   ownerID,
		timestamp: time.Now().UTC(),
	}
}

func (e *BucketCreatedEvent) EventType() string {
	return "bucket.created"
}

func (e *BucketCreatedEvent) BucketID() string {
	return e.bucketID
}

func (e *BucketCreatedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *BucketCreatedEvent) Name() string {
	return e.name
}

func (e *BucketCreatedEvent) Key() string {
	return e.key
}

func (e *BucketCreatedEvent) OwnerID() string {
	return e.ownerID
}
