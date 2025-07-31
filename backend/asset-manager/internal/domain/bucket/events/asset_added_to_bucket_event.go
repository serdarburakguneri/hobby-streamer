package events

import (
	"time"
)

type AssetAddedToBucketEvent struct {
	bucketID  string
	assetID   string
	timestamp time.Time
}

func NewAssetAddedToBucketEvent(bucketID, assetID string) *AssetAddedToBucketEvent {
	return &AssetAddedToBucketEvent{
		bucketID:  bucketID,
		assetID:   assetID,
		timestamp: time.Now().UTC(),
	}
}

func (e *AssetAddedToBucketEvent) EventType() string {
	return "asset.added_to_bucket"
}

func (e *AssetAddedToBucketEvent) BucketID() string {
	return e.bucketID
}

func (e *AssetAddedToBucketEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *AssetAddedToBucketEvent) AssetID() string {
	return e.assetID
}
