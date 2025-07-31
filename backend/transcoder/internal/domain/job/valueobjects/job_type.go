package valueobjects

type JobType string

const (
	JobTypeAnalyze   JobType = "analyze"
	JobTypeTranscode JobType = "transcode"
)

func (jt JobType) String() string {
	return string(jt)
}

func (jt JobType) IsAnalyze() bool {
	return jt == JobTypeAnalyze
}

func (jt JobType) IsTranscode() bool {
	return jt == JobTypeTranscode
}
