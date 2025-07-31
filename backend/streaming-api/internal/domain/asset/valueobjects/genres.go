package valueobjects

type Genres struct {
	values []Genre
}

func NewGenres(genres []string) (*Genres, error) {
	genreObjects := make([]Genre, 0, len(genres))
	for _, genreStr := range genres {
		genre, err := NewGenre(genreStr)
		if err != nil {
			return nil, err
		}
		genreObjects = append(genreObjects, *genre)
	}
	return &Genres{values: genreObjects}, nil
}

func (g Genres) Values() []Genre {
	return g.values
}

func (g Genres) Contains(genre Genre) bool {
	for _, existingGenre := range g.values {
		if existingGenre.Equals(genre) {
			return true
		}
	}
	return false
}
