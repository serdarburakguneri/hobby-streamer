package events

import (
	"time"
)

type AssetPublishedEvent struct {
	assetID   string
	publishAt time.Time
	regions   []string
	timestamp time.Time
}

func NewAssetPublishedEvent(assetID string, publishAt time.Time, regions []string) *AssetPublishedEvent {
	return &AssetPublishedEvent{
		assetID:   assetID,
		publishAt: publishAt,
		regions:   regions,
		timestamp: time.Now().UTC(),
	}
}

func (e *AssetPublishedEvent) EventType() string {
	return "asset.published"
}

func (e *AssetPublishedEvent) AssetID() string {
	return e.assetID
}

func (e *AssetPublishedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *AssetPublishedEvent) PublishAt() time.Time {
	return e.publishAt
}

func (e *AssetPublishedEvent) Regions() []string {
	return e.regions
}
