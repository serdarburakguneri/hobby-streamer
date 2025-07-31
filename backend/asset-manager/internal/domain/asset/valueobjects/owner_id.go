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
		return nil, errors.New("owner ID cannot be empty")
	}

	if len(value) > 100 {
		return nil, errors.New("owner ID too long")
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, value)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("owner ID contains invalid characters")
	}

	return &OwnerID{value: value}, nil
}

func (o OwnerID) Value() string {
	return o.value
}

func (o OwnerID) Equals(other OwnerID) bool {
	return o.value == other.value
}
