package valueobjects

import (
	"errors"
)

type VideoType string

const (
	VideoTypeMain    VideoType = "main"
	VideoTypeTrailer VideoType = "trailer"
	VideoTypeTeaser  VideoType = "teaser"
	VideoTypeBehind  VideoType = "behind"
	VideoTypeExtra   VideoType = "extra"
)

func NewVideoType(value string) (*VideoType, error) {
	if value == "" {
		return nil, errors.New("video type cannot be empty")
	}

	validTypes := []VideoType{
		VideoTypeMain,
		VideoTypeTrailer,
		VideoTypeTeaser,
		VideoTypeBehind,
		VideoTypeExtra,
	}

	for _, vType := range validTypes {
		if VideoType(value) == vType {
			vt := VideoType(value)
			return &vt, nil
		}
	}

	return nil, errors.New("invalid video type")
}

func (vt VideoType) Value() string {
	return string(vt)
}

func (vt VideoType) Equals(other VideoType) bool {
	return vt == other
}

func (vt VideoType) IsMain() bool {
	return vt == VideoTypeMain
}
