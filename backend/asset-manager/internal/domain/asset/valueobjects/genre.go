package valueobjects

import (
	"errors"
)

type Genre struct {
	value string
}

func NewGenre(value string) (*Genre, error) {
	if value == "" {
		return nil, errors.New("genre cannot be empty")
	}

	if len(value) > 50 {
		return nil, errors.New("genre too long")
	}

	validGenres := []string{
		"action", "adventure", "animation", "comedy", "crime", "documentary",
		"drama", "family", "fantasy", "horror", "mystery", "romance",
		"science-fiction", "thriller", "war", "western",
	}
	isValid := false
	for _, validGenre := range validGenres {
		if value == validGenre {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("invalid genre")
	}

	return &Genre{value: value}, nil
}

func (g Genre) Value() string {
	return g.value
}

func (g Genre) Equals(other Genre) bool {
	return g.value == other.value
}
