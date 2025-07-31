package valueobjects

import (
	"errors"
)

type VideoFormat string

const (
	VideoFormatRaw  VideoFormat = "raw"
	VideoFormatHLS  VideoFormat = "hls"
	VideoFormatDASH VideoFormat = "dash"
	VideoFormatMP4  VideoFormat = "mp4"
	VideoFormatWebM VideoFormat = "webm"
)

func NewVideoFormat(value string) (*VideoFormat, error) {
	if value == "" {
		return nil, errors.New("video format cannot be empty")
	}

	validFormats := []VideoFormat{
		VideoFormatRaw,
		VideoFormatHLS,
		VideoFormatDASH,
		VideoFormatMP4,
		VideoFormatWebM,
	}

	for _, format := range validFormats {
		if VideoFormat(value) == format {
			vf := VideoFormat(value)
			return &vf, nil
		}
	}

	return nil, errors.New("invalid video format")
}

func (vf VideoFormat) Value() string {
	return string(vf)
}

func (vf VideoFormat) Equals(other VideoFormat) bool {
	return vf == other
}

func (vf VideoFormat) IsStreaming() bool {
	return vf == VideoFormatHLS || vf == VideoFormatDASH
}

func (vf VideoFormat) IsDownloadable() bool {
	return vf == VideoFormatMP4 || vf == VideoFormatWebM
}
