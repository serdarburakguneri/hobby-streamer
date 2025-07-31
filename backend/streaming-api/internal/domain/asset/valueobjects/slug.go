package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type Slug struct {
	value string
}

func NewSlug(value string) (*Slug, error) {
	if value == "" {
		return nil, ErrInvalidSlug
	}

	if len(value) < 3 || len(value) > 100 {
		return nil, ErrInvalidSlug
	}

	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(value) {
		return nil, ErrInvalidSlug
	}

	return &Slug{value: strings.ToLower(value)}, nil
}

func (s Slug) Value() string {
	return s.value
}

func (s Slug) Equals(other Slug) bool {
	return s.value == other.value
}

var ErrInvalidSlug = errors.New("invalid slug")
