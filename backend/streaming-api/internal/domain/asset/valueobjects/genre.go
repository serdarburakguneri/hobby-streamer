package valueobjects

import (
	"errors"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type Genre struct {
	value string
}

var allowedGenres = constants.AllowedGenres

func NewGenre(value string) (*Genre, error) {
	if value == "" {
		return nil, ErrInvalidGenre
	}
	if _, ok := allowedGenres[value]; !ok {
		return nil, ErrInvalidGenre
	}
	return &Genre{value: strings.TrimSpace(value)}, nil
}

func (g Genre) Value() string {
	return g.value
}

func (g Genre) Equals(other Genre) bool {
	return g.value == other.value
}

var ErrInvalidGenre = errors.New("invalid genre")
