package valueobjects

import "fmt"

type ID struct {
	value string
}

func NewID(value string) (*ID, error) {
	if value == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	return &ID{value: value}, nil
}

func (id ID) Value() string {
	return id.value
}

func (id ID) String() string {
	return id.value
}
