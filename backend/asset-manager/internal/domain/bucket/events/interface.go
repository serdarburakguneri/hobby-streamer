package events

import (
	"time"
)

type BucketEvent interface {
	EventType() string
	BucketID() string
	Timestamp() time.Time
}
