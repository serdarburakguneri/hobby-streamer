package valueobjects

import (
	"errors"
	"strings"
)

type ValidatedString struct {
	value     string
	maxLength int
	fieldName string
}

func NewValidatedString(value string, maxLength int, fieldName string) (*ValidatedString, error) {
	if value == "" {
		return nil, errors.New(fieldName + " cannot be empty")
	}

	if len(value) > maxLength {
		return nil, errors.New(fieldName + " too long")
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, errors.New(fieldName + " cannot be empty after trimming")
	}

	return &ValidatedString{
		value:     trimmed,
		maxLength: maxLength,
		fieldName: fieldName,
	}, nil
}

func (vs ValidatedString) Value() string {
	return vs.value
}

func (vs ValidatedString) Equals(other ValidatedString) bool {
	return vs.value == other.value
}

func (vs ValidatedString) MaxLength() int {
	return vs.maxLength
}

func (vs ValidatedString) FieldName() string {
	return vs.fieldName
}
