package events

type CompletedEvent interface {
	Type() string
	Data() interface{}
	Topic() string
	CloudEventType() string
	ID() string
}

type JobCompletedBase struct {
	JobID        string `json:"jobId"`
	AssetID      string `json:"assetId"`
	VideoID      string `json:"videoId"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	CompletedAt  string `json:"completedAt"`
}

func (b JobCompletedBase) ID() string { return b.JobID }
