package valueobjects

import (
	"errors"
)

type VideoStatus string

const (
	VideoStatusPending     VideoStatus = "pending"
	VideoStatusAnalyzing   VideoStatus = "analyzing"
	VideoStatusTranscoding VideoStatus = "transcoding"
	VideoStatusReady       VideoStatus = "ready"
	VideoStatusFailed      VideoStatus = "failed"
)

func NewVideoStatus(value string) (*VideoStatus, error) {
	if value == "" {
		return nil, errors.New("video status cannot be empty")
	}

	validStatuses := []VideoStatus{
		VideoStatusPending,
		VideoStatusAnalyzing,
		VideoStatusTranscoding,
		VideoStatusReady,
		VideoStatusFailed,
	}

	for _, status := range validStatuses {
		if VideoStatus(value) == status {
			vs := VideoStatus(value)
			return &vs, nil
		}
	}

	return nil, errors.New("invalid video status")
}

func (vs VideoStatus) Value() string {
	return string(vs)
}

func (vs VideoStatus) Equals(other VideoStatus) bool {
	return vs == other
}

func (vs VideoStatus) IsReady() bool {
	return vs == VideoStatusReady
}

func (vs VideoStatus) IsFailed() bool {
	return vs == VideoStatusFailed
}

func (vs VideoStatus) IsProcessing() bool {
	return vs == VideoStatusPending || vs == VideoStatusAnalyzing || vs == VideoStatusTranscoding
}
