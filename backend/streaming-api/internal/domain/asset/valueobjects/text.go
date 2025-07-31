package valueobjects

import (
	"errors"
	"strings"
)

type Text struct {
	value string
	kind  string
}

func NewText(value, kind string, maxLength int) (*Text, error) {
	if value == "" {
		return nil, ErrInvalidText
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, ErrInvalidText
	}

	if maxLength > 0 && len(trimmed) > maxLength {
		return nil, ErrInvalidText
	}

	return &Text{value: trimmed, kind: kind}, nil
}

func (t Text) Value() string {
	return t.value
}

func (t Text) Kind() string {
	return t.kind
}

func (t Text) Equals(other Text) bool {
	return t.value == other.value && t.kind == other.kind
}

var ErrInvalidText = errors.New("invalid text")
