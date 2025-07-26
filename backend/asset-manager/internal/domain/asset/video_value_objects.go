package asset

import (
	"regexp"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

var (
	ErrInvalidVideoID          = pkgerrors.NewValidationError("invalid video id", nil)
	ErrInvalidVideoLabel       = pkgerrors.NewValidationError("invalid video label", nil)
	ErrInvalidVideoContentType = pkgerrors.NewValidationError("invalid video content type", nil)
)

type VideoID struct {
	value string
}

func NewVideoID(value string) (*VideoID, error) {
	if value == "" {
		return nil, ErrInvalidVideoID
	}
	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidVideoID
	}
	return &VideoID{value: value}, nil
}

func (id VideoID) Value() string {
	return id.value
}

func (id VideoID) Equals(other VideoID) bool {
	return id.value == other.value
}

type VideoLabel struct {
	value string
}

func NewVideoLabel(value string) (*VideoLabel, error) {
	if value == "" {
		return nil, ErrInvalidVideoLabel
	}
	if len(value) > 100 {
		return nil, ErrInvalidVideoLabel
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, ErrInvalidVideoLabel
	}
	return &VideoLabel{value: trimmed}, nil
}

func (l VideoLabel) Value() string {
	return l.value
}

func (l VideoLabel) Equals(other VideoLabel) bool {
	return l.value == other.value
}

type VideoContentType struct {
	value string
}

func NewVideoContentType(value string) (*VideoContentType, error) {
	if value == "" {
		return &VideoContentType{value: ""}, nil
	}
	if len(value) > 100 {
		return nil, ErrInvalidVideoContentType
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return &VideoContentType{value: ""}, nil
	}
	return &VideoContentType{value: trimmed}, nil
}

func (c VideoContentType) Value() string {
	return c.value
}

func (c VideoContentType) Equals(other VideoContentType) bool {
	return c.value == other.value
}
