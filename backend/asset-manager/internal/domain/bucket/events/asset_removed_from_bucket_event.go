package events

import (
	"time"
)

type AssetRemovedFromBucketEvent struct {
	bucketID  string
	assetID   string
	timestamp time.Time
}

func NewAssetRemovedFromBucketEvent(bucketID, assetID string) *AssetRemovedFromBucketEvent {
	return &AssetRemovedFromBucketEvent{
		bucketID:  bucketID,
		assetID:   assetID,
		timestamp: time.Now().UTC(),
	}
}

func (e *AssetRemovedFromBucketEvent) EventType() string {
	return "asset.removed_from_bucket"
}

func (e *AssetRemovedFromBucketEvent) BucketID() string {
	return e.bucketID
}

func (e *AssetRemovedFromBucketEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *AssetRemovedFromBucketEvent) AssetID() string {
	return e.assetID
}
