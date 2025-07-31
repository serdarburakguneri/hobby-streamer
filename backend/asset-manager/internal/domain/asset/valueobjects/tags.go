package valueobjects

import (
	"errors"
)

type Tags struct {
	values []Tag
}

func NewTags(values []string) (*Tags, error) {
	if len(values) > 20 {
		return nil, errors.New("too many tags")
	}

	tags := make([]Tag, 0, len(values))
	for _, value := range values {
		tag, err := NewTag(value)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	return &Tags{values: tags}, nil
}

func (t Tags) Values() []Tag {
	return t.values
}

func (t Tags) Count() int {
	return len(t.values)
}

func (t Tags) Contains(tag Tag) bool {
	for _, t := range t.values {
		if t.Equals(tag) {
			return true
		}
	}
	return false
}
