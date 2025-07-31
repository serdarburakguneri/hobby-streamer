package valueobjects

type TranscodingInfo struct {
	width       int
	height      int
	duration    float64
	bitrate     int
	codec       string
	size        int64
	contentType ContentType
}

func NewTranscodingInfo(width, height int, duration float64, bitrate int, codec string, size int64, contentType ContentType) *TranscodingInfo {
	return &TranscodingInfo{
		width:       width,
		height:      height,
		duration:    duration,
		bitrate:     bitrate,
		codec:       codec,
		size:        size,
		contentType: contentType,
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
