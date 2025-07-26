package constants

const (
	VideoStreamingFormatRaw  = "raw"
	VideoStreamingFormatHLS  = "hls"
	VideoStreamingFormatDASH = "dash"
)

var AllowedVideoStreamingFormats = map[string]struct{}{
	VideoStreamingFormatRaw:  {},
	VideoStreamingFormatHLS:  {},
	VideoStreamingFormatDASH: {},
}

func IsValidVideoStreamingFormat(f string) bool {
	_, ok := AllowedVideoStreamingFormats[f]
	return ok
}
