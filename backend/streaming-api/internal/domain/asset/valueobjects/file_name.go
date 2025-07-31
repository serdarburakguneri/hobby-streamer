package valueobjects

import (
	"errors"
	"regexp"
)

type FileName struct {
	value string
}

func NewFileName(value string) (*FileName, error) {
	if value == "" {
		return nil, ErrInvalidFileName
	}

	if len(value) > 255 {
		return nil, ErrInvalidFileName
	}

	fileNameRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !fileNameRegex.MatchString(value) {
		return nil, ErrInvalidFileName
	}

	return &FileName{value: value}, nil
}

func (f FileName) Value() string {
	return f.value
}

func (f FileName) Equals(other FileName) bool {
	return f.value == other.value
}

var ErrInvalidFileName = errors.New("invalid file name")
