package constants

const (
	VideoFormatHLS  = "hls"
	VideoFormatDASH = "dash"
	VideoFormatRAW  = "raw"
	VideoFormatMP4  = "mp4"
	VideoFormatWEBM = "webm"
	VideoFormatAVI  = "avi"
	VideoFormatMOV  = "mov"
	VideoFormatMKV  = "mkv"
)

var AllowedVideoFormats = map[string]struct{}{
	VideoFormatHLS:  {},
	VideoFormatDASH: {},
	VideoFormatRAW:  {},
	VideoFormatMP4:  {},
	VideoFormatWEBM: {},
	VideoFormatAVI:  {},
	VideoFormatMOV:  {},
	VideoFormatMKV:  {},
}

func IsValidVideoFormat(f string) bool {
	_, ok := AllowedVideoFormats[f]
	return ok
}
