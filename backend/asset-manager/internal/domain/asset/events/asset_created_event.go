package events

import (
	"time"
)

type AssetCreatedEvent struct {
	assetID   string
	slug      string
	title     string
	assetType string
	ownerID   string
	timestamp time.Time
}

func NewAssetCreatedEvent(assetID, slug, title, assetType, ownerID string) *AssetCreatedEvent {
	return &AssetCreatedEvent{
		assetID:   assetID,
		slug:      slug,
		title:     title,
		assetType: assetType,
		ownerID:   ownerID,
		timestamp: time.Now().UTC(),
	}
}

func (e *AssetCreatedEvent) EventType() string {
	return "asset.created"
}

func (e *AssetCreatedEvent) AssetID() string {
	return e.assetID
}

func (e *AssetCreatedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *AssetCreatedEvent) Slug() string {
	return e.slug
}

func (e *AssetCreatedEvent) Title() string {
	return e.title
}

func (e *AssetCreatedEvent) AssetType() string {
	return e.assetType
}

func (e *AssetCreatedEvent) OwnerID() string {
	return e.ownerID
}
