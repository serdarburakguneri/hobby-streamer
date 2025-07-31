package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type Tags struct {
	values []string
}

func NewTags(tags []string) (*Tags, error) {
	if len(tags) > 20 {
		return nil, ErrTooManyTags
	}

	validatedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag) > 30 {
			return nil, ErrInvalidTag
		}

		tagRegex := regexp.MustCompile(`^[a-zA-Z0-9\s-]+$`)
		if !tagRegex.MatchString(tag) {
			return nil, ErrInvalidTag
		}

		validatedTags = append(validatedTags, strings.TrimSpace(tag))
	}

	return &Tags{values: validatedTags}, nil
}

func (t Tags) Values() []string {
	return t.values
}

func (t Tags) Contains(tag string) bool {
	for _, existingTag := range t.values {
		if existingTag == tag {
			return true
		}
	}
	return false
}

var (
	ErrInvalidTag  = errors.New("invalid tag")
	ErrTooManyTags = errors.New("too many tags")
)
