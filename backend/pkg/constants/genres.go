package constants

const (
	GenreAction      = "action"
	GenreDrama       = "drama"
	GenreComedy      = "comedy"
	GenreHorror      = "horror"
	GenreSciFi       = "sci_fi"
	GenreRomance     = "romance"
	GenreThriller    = "thriller"
	GenreFantasy     = "fantasy"
	GenreDocumentary = "documentary"
	GenreMusic       = "music"
	GenreNews        = "news"
	GenreSports      = "sports"
	GenreKids        = "kids"
	GenreEducational = "educational"
	GenreWestern     = "western"
	GenreAnimation   = "animation"
	GenreFamily      = "family"
	GenreMystery     = "mystery"
	GenreWar         = "war"
	GenreCrime       = "crime"
)

var AllowedGenres = map[string]struct{}{
	GenreAction:      {},
	GenreDrama:       {},
	GenreComedy:      {},
	GenreHorror:      {},
	GenreSciFi:       {},
	GenreRomance:     {},
	GenreThriller:    {},
	GenreFantasy:     {},
	GenreDocumentary: {},
	GenreMusic:       {},
	GenreNews:        {},
	GenreSports:      {},
	GenreKids:        {},
	GenreEducational: {},
	GenreWestern:     {},
	GenreAnimation:   {},
	GenreFamily:      {},
	GenreMystery:     {},
	GenreWar:         {},
	GenreCrime:       {},
}

func IsValidGenre(g string) bool {
	_, ok := AllowedGenres[g]
	return ok
}
