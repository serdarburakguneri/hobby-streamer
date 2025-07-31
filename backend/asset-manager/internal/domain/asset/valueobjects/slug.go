package valueobjects

import (
	"regexp"
	"strings"
)

type Slug struct {
	value string
}

func NewSlug(value string) (*Slug, error) {
	validatedString, err := NewValidatedString(value, 100, "slug")
	if err != nil {
		return nil, err
	}
	return &Slug{value: validatedString.Value()}, nil
}

func (s Slug) Value() string {
	return s.value
}

func (s Slug) Equals(other Slug) bool {
	return s.value == other.value
}

func GenerateSlug(title string) (string, error) {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	_, err := NewSlug(slug)
	if err != nil {
		return "", err
	}

	return slug, nil
}
