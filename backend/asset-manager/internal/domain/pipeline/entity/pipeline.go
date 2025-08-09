package entity

import (
	"time"
)

type StepState struct {
	Status        string     `json:"status"`
	StartedAt     time.Time  `json:"startedAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	ErrorMessage  string     `json:"errorMessage,omitempty"`
	JobID         string     `json:"jobId,omitempty"`
	CorrelationID string     `json:"correlationId,omitempty"`
}

type Pipeline struct {
	AssetID   string               `json:"assetId"`
	VideoID   string               `json:"videoId"`
	Steps     map[string]StepState `json:"steps"`
	UpdatedAt time.Time            `json:"updatedAt"`
	CreatedAt time.Time            `json:"createdAt"`
}

func NewPipeline(assetID, videoID string) *Pipeline {
	now := time.Now().UTC()
	return &Pipeline{AssetID: assetID, VideoID: videoID, Steps: map[string]StepState{}, UpdatedAt: now, CreatedAt: now}
}

func (p *Pipeline) SetRequested(step, jobID, correlationID string) {
	p.Steps[step] = StepState{Status: "requested", StartedAt: time.Now().UTC(), JobID: jobID, CorrelationID: correlationID}
	p.UpdatedAt = time.Now().UTC()
}

func (p *Pipeline) SetCompleted(step string) {
	now := time.Now().UTC()
	s := p.Steps[step]
	s.Status = "completed"
	s.CompletedAt = &now
	p.Steps[step] = s
	p.UpdatedAt = now
}

func (p *Pipeline) SetFailed(step, errMsg string) {
	now := time.Now().UTC()
	s := p.Steps[step]
	s.Status = "failed"
	s.ErrorMessage = errMsg
	s.CompletedAt = &now
	p.Steps[step] = s
	p.UpdatedAt = now
}
