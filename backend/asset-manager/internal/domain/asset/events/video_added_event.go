package events

import (
	"time"
)

type VideoAddedEvent struct {
	assetID   string
	videoID   string
	label     string
	format    string
	timestamp time.Time
}

func NewVideoAddedEvent(assetID, videoID, label, format string) *VideoAddedEvent {
	return &VideoAddedEvent{
		assetID:   assetID,
		videoID:   videoID,
		label:     label,
		format:    format,
		timestamp: time.Now().UTC(),
	}
}

func (e *VideoAddedEvent) EventType() string {
	return "video.added"
}

func (e *VideoAddedEvent) AssetID() string {
	return e.assetID
}

func (e *VideoAddedEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *VideoAddedEvent) VideoID() string {
	return e.videoID
}

func (e *VideoAddedEvent) Label() string {
	return e.label
}

func (e *VideoAddedEvent) Format() string {
	return e.format
}
