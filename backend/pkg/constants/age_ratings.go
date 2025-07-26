package constants

const (
	AgeRatingG    = "G"
	AgeRatingPG   = "PG"
	AgeRatingPG13 = "PG-13"
	AgeRatingR    = "R"
	AgeRatingNC17 = "NC-17"
	AgeRatingTVY  = "TV-Y"
	AgeRatingTVY7 = "TV-Y7"
	AgeRatingTVG  = "TV-G"
	AgeRatingTVPG = "TV-PG"
	AgeRatingTV14 = "TV-14"
	AgeRatingTVMA = "TV-MA"
)

var AllowedAgeRatings = map[string]struct{}{
	AgeRatingG:    {},
	AgeRatingPG:   {},
	AgeRatingPG13: {},
	AgeRatingR:    {},
	AgeRatingNC17: {},
	AgeRatingTVY:  {},
	AgeRatingTVY7: {},
	AgeRatingTVG:  {},
	AgeRatingTVPG: {},
	AgeRatingTV14: {},
	AgeRatingTVMA: {},
}

func IsValidAgeRating(r string) bool {
	_, ok := AllowedAgeRatings[r]
	return ok
}
