package constants

const (
	VideoTypeMain      = "main"
	VideoTypeTrailer   = "trailer"
	VideoTypeBehind    = "behind"
	VideoTypeInterview = "interview"
)

var AllowedVideoTypes = map[string]struct{}{
	VideoTypeMain:      {},
	VideoTypeTrailer:   {},
	VideoTypeBehind:    {},
	VideoTypeInterview: {},
}

func IsValidVideoType(t string) bool {
	_, ok := AllowedVideoTypes[t]
	return ok
}
