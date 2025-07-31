package valueobjects

type JobFormat string

const (
	JobFormatHLS  JobFormat = "hls"
	JobFormatDASH JobFormat = "dash"
)

func (jf JobFormat) String() string {
	return string(jf)
}

func (jf JobFormat) IsHLS() bool {
	return jf == JobFormatHLS
}

func (jf JobFormat) IsDASH() bool {
	return jf == JobFormatDASH
}
