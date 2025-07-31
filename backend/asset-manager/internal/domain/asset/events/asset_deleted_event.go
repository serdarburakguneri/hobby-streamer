package events

import (
	"time"
)

type AssetDeletedEvent struct {
	assetID   string
	timestamp time.Time
}

func NewAssetDeletedEvent(assetID string) *AssetDeletedEvent {
	return &AssetDeletedEvent{
		assetID:   assetID,
		timestamp: time.Now().UTC(),
	}
}

func (e *AssetDeletedEvent) EventType() string {
	return "asset.deleted"
}

func (e *AssetDeletedEvent) AssetID() string {
	return e.assetID
}

func (e *AssetDeletedEvent) Timestamp() time.Time {
	return e.timestamp
}
