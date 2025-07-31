package valueobjects

import (
	"errors"
	"regexp"
)

type ID struct {
	value     string
	idType    string
	maxLength int
}

func NewID(value, idType string, maxLength int) (*ID, error) {
	if value == "" {
		return nil, errors.New(idType + " ID cannot be empty")
	}

	if len(value) > maxLength {
		return nil, errors.New(idType + " ID too long")
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, value)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New(idType + " ID contains invalid characters")
	}

	return &ID{
		value:     value,
		idType:    idType,
		maxLength: maxLength,
	}, nil
}

func (id ID) Value() string {
	return id.value
}

func (id ID) Type() string {
	return id.idType
}

func (id ID) Equals(other ID) bool {
	return id.value == other.value && id.idType == other.idType
}
