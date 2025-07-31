package valueobjects

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

func (js JobStatus) String() string {
	return string(js)
}

func (js JobStatus) IsPending() bool {
	return js == JobStatusPending
}

func (js JobStatus) IsRunning() bool {
	return js == JobStatusRunning
}

func (js JobStatus) IsCompleted() bool {
	return js == JobStatusCompleted
}

func (js JobStatus) IsFailed() bool {
	return js == JobStatusFailed
}
