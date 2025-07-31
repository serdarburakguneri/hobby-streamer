package valueobjects

import (
	"errors"
	"strings"
)

type Description struct {
	value string
}

func NewDescription(value string) (*Description, error) {
	if len(value) > 1000 {
		return nil, ErrInvalidDescription
	}

	return &Description{value: strings.TrimSpace(value)}, nil
}

func (d Description) Value() string {
	return d.value
}

func (d Description) Equals(other Description) bool {
	return d.value == other.value
}

var ErrInvalidDescription = errors.New("invalid description")
