package valueobjects

import (
	"errors"
)

type Genres struct {
	values []Genre
}

func NewGenres(values []string) (*Genres, error) {
	if len(values) > 10 {
		return nil, errors.New("too many genres")
	}

	genres := make([]Genre, 0, len(values))
	for _, value := range values {
		genre, err := NewGenre(value)
		if err != nil {
			return nil, err
		}
		genres = append(genres, *genre)
	}

	return &Genres{values: genres}, nil
}

func (g Genres) Values() []Genre {
	return g.values
}

func (g Genres) Count() int {
	return len(g.values)
}

func (g Genres) Contains(genre Genre) bool {
	for _, g := range g.values {
		if g.Equals(genre) {
			return true
		}
	}
	return false
}
