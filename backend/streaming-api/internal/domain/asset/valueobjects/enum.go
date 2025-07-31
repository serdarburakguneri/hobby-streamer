package valueobjects

import (
	"errors"
)

type Enum struct {
	value string
	kind  string
}

func NewEnum(value, kind string, allowedValues map[string]bool) (*Enum, error) {
	if value == "" {
		return nil, ErrInvalidEnum
	}

	if allowedValues != nil && !allowedValues[value] {
		return nil, ErrInvalidEnum
	}

	return &Enum{value: value, kind: kind}, nil
}

func (e Enum) Value() string {
	return e.value
}

func (e Enum) Kind() string {
	return e.kind
}

func (e Enum) Equals(other Enum) bool {
	return e.value == other.value && e.kind == other.kind
}

var ErrInvalidEnum = errors.New("invalid enum value")
