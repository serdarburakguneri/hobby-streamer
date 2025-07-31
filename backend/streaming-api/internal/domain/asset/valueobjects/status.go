package valueobjects

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type Status struct {
	value string
}

func NewStatus(value string) (*Status, error) {
	if value == "" {
		return nil, ErrInvalidStatus
	}

	validStatuses := map[string]bool{
		constants.AssetStatusDraft:     true,
		"processing":                   true,
		constants.VideoStatusReady:     true,
		constants.AssetStatusPublished: true,
		"archived":                     true,
		"deleted":                      true,
	}

	if !validStatuses[value] {
		return nil, ErrInvalidStatus
	}

	return &Status{value: value}, nil
}

func (s Status) Value() string {
	return s.value
}

func (s Status) Equals(other Status) bool {
	return s.value == other.value
}

var ErrInvalidStatus = errors.New("invalid status")
