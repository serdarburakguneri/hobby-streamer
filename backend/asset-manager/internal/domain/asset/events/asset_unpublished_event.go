package events

import (
	"time"
)

type AssetUnpublishedEvent struct {
	assetID   string
	timestamp time.Time
}

func NewAssetUnpublishedEvent(assetID string) *AssetUnpublishedEvent {
	return &AssetUnpublishedEvent{
		assetID:   assetID,
		timestamp: time.Now().UTC(),
	}
}

func (e *AssetUnpublishedEvent) EventType() string {
	return "asset.unpublished"
}

func (e *AssetUnpublishedEvent) AssetID() string {
	return e.assetID
}

func (e *AssetUnpublishedEvent) Timestamp() time.Time {
	return e.timestamp
}
