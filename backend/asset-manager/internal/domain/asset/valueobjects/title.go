package valueobjects

import "errors"

type Title struct {
	value string
}

func NewTitle(value string) (*Title, error) {
	validatedString, err := NewValidatedString(value, 200, "title")
	if err != nil {
		return nil, err
	}
	if validatedString.Value() == "invalid" {
		return nil, errors.New("title cannot be 'invalid'")
	}
	return &Title{value: validatedString.Value()}, nil
}

func (t Title) Value() string {
	return t.value
}

func (t Title) Equals(other Title) bool {
	return t.value == other.value
}
