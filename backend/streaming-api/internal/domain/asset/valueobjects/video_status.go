package valueobjects

import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"

import (
	"errors"
)

type VideoStatus string

func NewVideoStatus(value string) (*VideoStatus, error) {
	if value == "" {
		return nil, errors.New("video status cannot be empty")
	}

	if !constants.IsValidVideoStatus(value) {
		return nil, errors.New("invalid video status")
	}

	vs := VideoStatus(value)
	return &vs, nil
}

func (vs VideoStatus) Value() string {
	return string(vs)
}

func (vs VideoStatus) Equals(other VideoStatus) bool {
	return vs == other
}

func (vs VideoStatus) IsReady() bool {
	return vs == VideoStatus(constants.VideoStatusReady)
}

func (vs VideoStatus) IsFailed() bool {
	return vs == VideoStatus(constants.VideoStatusFailed)
}

func (vs VideoStatus) IsProcessing() bool {
	return vs == VideoStatus(constants.VideoStatusPending) ||
		vs == VideoStatus(constants.VideoStatusAnalyzing) ||
		vs == VideoStatus(constants.VideoStatusTranscoding)
}
