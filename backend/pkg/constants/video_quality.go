package constants

const (
	VideoQualityMain = "main"
	VideoQualitySD   = "sd"
	VideoQualityHD   = "hd"
	VideoQualityFHD  = "fhd"
	VideoQuality4K   = "4k"
)

var AllowedVideoQualities = map[string]struct{}{
	VideoQualityMain: {},
	VideoQualitySD:   {},
	VideoQualityHD:   {},
	VideoQualityFHD:  {},
	VideoQuality4K:   {},
}

func IsValidVideoQuality(q string) bool {
	_, ok := AllowedVideoQualities[q]
	return ok
}
