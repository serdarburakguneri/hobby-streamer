package valueobjects

import (
	"errors"
	"strings"
)

type Title struct {
	value string
}

func NewTitle(value string) (*Title, error) {
	if value == "" {
		return nil, ErrInvalidTitle
	}

	if len(value) > 255 {
		return nil, ErrInvalidTitle
	}

	return &Title{value: strings.TrimSpace(value)}, nil
}

func (t Title) Value() string {
	return t.value
}

func (t Title) Equals(other Title) bool {
	return t.value == other.value
}

var ErrInvalidTitle = errors.New("invalid title")
