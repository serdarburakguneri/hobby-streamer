package valueobjects

import (
	"errors"
	"regexp"
)

type OwnerID struct {
	value string
}

func NewOwnerID(value string) (*OwnerID, error) {
	if value == "" {
		return nil, ErrInvalidOwnerID
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(value) {
		return nil, ErrInvalidOwnerID
	}

	return &OwnerID{value: value}, nil
}

func (o OwnerID) Value() string {
	return o.value
}

func (o OwnerID) Equals(other OwnerID) bool {
	return o.value == other.value
}

var ErrInvalidOwnerID = errors.New("invalid owner ID")
