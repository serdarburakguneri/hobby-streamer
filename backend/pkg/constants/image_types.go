package constants

const (
	ImageTypePoster       = "poster"
	ImageTypeBackdrop     = "backdrop"
	ImageTypeThumbnail    = "thumbnail"
	ImageTypeLogo         = "logo"
	ImageTypeBanner       = "banner"
	ImageTypeHero         = "hero"
	ImageTypeScreenshot   = "screenshot"
	ImageTypeBehindScenes = "behind_scenes"
	ImageTypeInterview    = "interview"
)

var AllowedImageTypes = map[string]struct{}{
	ImageTypePoster:       {},
	ImageTypeBackdrop:     {},
	ImageTypeThumbnail:    {},
	ImageTypeLogo:         {},
	ImageTypeBanner:       {},
	ImageTypeHero:         {},
	ImageTypeScreenshot:   {},
	ImageTypeBehindScenes: {},
	ImageTypeInterview:    {},
}

func IsValidImageType(t string) bool {
	_, ok := AllowedImageTypes[t]
	return ok
}
