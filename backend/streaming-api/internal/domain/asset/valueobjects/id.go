package valueobjects

import (
	"errors"
	"regexp"
)

type ID struct {
	value string
	kind  string
}

func NewID(value, kind string) (*ID, error) {
	if value == "" {
		return nil, ErrInvalidID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidID
	}

	return &ID{value: value, kind: kind}, nil
}

func (id ID) Value() string {
	return id.value
}

func (id ID) Kind() string {
	return id.kind
}

func (id ID) Equals(other ID) bool {
	return id.value == other.value && id.kind == other.kind
}

var ErrInvalidID = errors.New("invalid ID")
