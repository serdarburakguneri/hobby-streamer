package events

import (
	"time"
)

type VideoRemovedEvent struct {
	assetID   string
	videoID   string
	timestamp time.Time
}

func NewVideoRemovedEvent(assetID, videoID string) *VideoRemovedEvent {
	return &VideoRemovedEvent{
		assetID:   assetID,
		videoID:   videoID,
		timestamp: time.Now().UTC(),
	}
}

func (e *VideoRemovedEvent) EventType() string {
	return "video.removed"
}

func (e *VideoRemovedEvent) AssetID() string {
	return e.assetID
}

func (e *VideoRemovedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *VideoRemovedEvent) VideoID() string {
	return e.videoID
}
