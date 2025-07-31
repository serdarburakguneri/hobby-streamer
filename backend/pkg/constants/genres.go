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
}

func IsValidGenre(g string) bool {
	_, ok := AllowedGenres[g]
	return ok
}
