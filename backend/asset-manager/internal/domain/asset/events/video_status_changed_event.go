package events

import (
	"time"
)

type VideoStatusChangedEvent struct {
	assetID   string
	videoID   string
	oldStatus string
	newStatus string
	timestamp time.Time
}

func NewVideoStatusChangedEvent(assetID, videoID, oldStatus, newStatus string) *VideoStatusChangedEvent {
	return &VideoStatusChangedEvent{
		assetID:   assetID,
		videoID:   videoID,
		oldStatus: oldStatus,
		newStatus: newStatus,
		timestamp: time.Now().UTC(),
	}
}

func (e *VideoStatusChangedEvent) EventType() string {
	return "video.status_changed"
}

func (e *VideoStatusChangedEvent) AssetID() string {
	return e.assetID
}

func (e *VideoStatusChangedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *VideoStatusChangedEvent) VideoID() string {
	return e.videoID
}

func (e *VideoStatusChangedEvent) OldStatus() string {
	return e.oldStatus
}

func (e *VideoStatusChangedEvent) NewStatus() string {
	return e.newStatus
}
