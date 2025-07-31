package events

import (
	"time"
)

type AssetEvent interface {
	EventType() string
	AssetID() string
	Timestamp() time.Time
}
