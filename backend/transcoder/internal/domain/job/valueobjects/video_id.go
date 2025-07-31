package valueobjects

import "fmt"

type VideoID struct {
	value string
}

func NewVideoID(value string) (*VideoID, error) {
	if value == "" {
		return nil, fmt.Errorf("video ID cannot be empty")
	}
	return &VideoID{value: value}, nil
}

func (v VideoID) Value() string {
	return v.value
}

func (v VideoID) String() string {
	return v.value
}
