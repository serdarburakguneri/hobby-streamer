package events

import (
	"time"
)

type AssetUpdatedEvent struct {
	assetID       string
	updatedFields map[string]interface{}
	timestamp     time.Time
}

func NewAssetUpdatedEvent(assetID string, updatedFields map[string]interface{}) *AssetUpdatedEvent {
	return &AssetUpdatedEvent{
		assetID:       assetID,
		updatedFields: updatedFields,
		timestamp:     time.Now().UTC(),
	}
}

func (e *AssetUpdatedEvent) EventType() string {
	return "asset.updated"
}

func (e *AssetUpdatedEvent) AssetID() string {
	return e.assetID
}

func (e *AssetUpdatedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *AssetUpdatedEvent) UpdatedFields() map[string]interface{} {
	return e.updatedFields
}
