package events

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
)

type JobCompletedEvent struct {
	JobID        string      `json:"jobId"`
	AssetID      string      `json:"assetId"`
	VideoID      string      `json:"videoId"`
	JobType      string      `json:"jobType"`
	Format       string      `json:"format,omitempty"`
	Success      bool        `json:"success"`
	Metadata     interface{} `json:"metadata,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	CompletedAt  time.Time   `json:"completedAt"`
}

func NewJobCompletedEvent(job *entity.Job, success bool, metadata interface{}, errorMessage string) *JobCompletedEvent {
	return &JobCompletedEvent{
		JobID:        job.ID().Value(),
		AssetID:      job.AssetID().Value(),
		VideoID:      job.VideoID().Value(),
		JobType:      job.Type().String(),
		Format:       job.Format().String(),
		Success:      success,
		Metadata:     metadata,
		ErrorMessage: errorMessage,
		CompletedAt:  time.Now().UTC(),
	}
}

func (e *JobCompletedEvent) Type() string {
	return "job.completed"
}

func (e *JobCompletedEvent) Data() interface{} {
	return e
}
