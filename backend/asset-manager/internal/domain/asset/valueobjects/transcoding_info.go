package valueobjects

type TranscodingInfo struct {
	width           int
	height          int
	duration        float64
	bitrate         int
	codec           string
	size            int64
	contentType     ContentType
	videoCodec      string
	audioCodec      string
	frameRate       string
	audioChannels   int
	audioSampleRate int
}

func NewMediaInfo(width, height int, duration float64, bitrate int, codec string, size int64, contentType ContentType, videoCodec, audioCodec, frameRate string, audioChannels, audioSampleRate int) *TranscodingInfo {
	return &TranscodingInfo{
		width:           width,
		height:          height,
		duration:        duration,
		bitrate:         bitrate,
		codec:           codec,
		size:            size,
		contentType:     contentType,
		videoCodec:      videoCodec,
		audioCodec:      audioCodec,
		frameRate:       frameRate,
		audioChannels:   audioChannels,
		audioSampleRate: audioSampleRate,
	}
}

func (ti TranscodingInfo) Width() int {
	return ti.width
}

func (ti TranscodingInfo) Height() int {
	return ti.height
}

func (ti TranscodingInfo) Duration() float64 {
	return ti.duration
}

func (ti TranscodingInfo) Bitrate() int {
	return ti.bitrate
}

func (ti TranscodingInfo) Codec() string {
	return ti.codec
}

func (ti TranscodingInfo) Size() int64 {
	return ti.size
}

func (ti TranscodingInfo) ContentType() ContentType {
	return ti.contentType
}

func (ti TranscodingInfo) VideoCodec() string   { return ti.videoCodec }
func (ti TranscodingInfo) AudioCodec() string   { return ti.audioCodec }
func (ti TranscodingInfo) FrameRate() string    { return ti.frameRate }
func (ti TranscodingInfo) AudioChannels() int   { return ti.audioChannels }
func (ti TranscodingInfo) AudioSampleRate() int { return ti.audioSampleRate }
