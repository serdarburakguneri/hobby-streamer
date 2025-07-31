package entity

import (
	"fmt"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

type Job struct {
	id          valueobjects.JobID
	jobType     valueobjects.JobType
	format      valueobjects.JobFormat
	assetID     valueobjects.AssetID
	videoID     valueobjects.VideoID
	input       string
	output      string
	quality     string
	status      valueobjects.JobStatus
	progress    float64
	error       string
	metadata    *valueobjects.VideoMetadata
	startedAt   *time.Time
	completedAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewAnalyzeJob(assetID valueobjects.AssetID, videoID valueobjects.VideoID, input string) *Job {
	now := time.Now().UTC()
	jid, _ := valueobjects.NewJobID(valueobjects.GenerateJobID())
	return &Job{
		id:        *jid,
		jobType:   valueobjects.JobTypeAnalyze,
		assetID:   assetID,
		videoID:   videoID,
		input:     input,
		status:    valueobjects.JobStatusPending,
		progress:  0.0,
		createdAt: now,
		updatedAt: now,
	}
}

func NewTranscodeJob(assetID valueobjects.AssetID, videoID valueobjects.VideoID, input, output, quality string, format valueobjects.JobFormat) *Job {
	now := time.Now().UTC()
	jid, _ := valueobjects.NewJobID(valueobjects.GenerateJobID())
	return &Job{
		id:        *jid,
		jobType:   valueobjects.JobTypeTranscode,
		format:    format,
		assetID:   assetID,
		videoID:   videoID,
		input:     input,
		output:    output,
		quality:   quality,
		status:    valueobjects.JobStatusPending,
		progress:  0.0,
		createdAt: now,
		updatedAt: now,
	}
}

func (j *Job) ID() valueobjects.JobID {
	return j.id
}

func (j *Job) Type() valueobjects.JobType {
	return j.jobType
}

func (j *Job) Format() valueobjects.JobFormat {
	return j.format
}

func (j *Job) AssetID() valueobjects.AssetID {
	return j.assetID
}

func (j *Job) VideoID() valueobjects.VideoID {
	return j.videoID
}

func (j *Job) Input() string {
	return j.input
}

func (j *Job) Output() string {
	return j.output
}

func (j *Job) Quality() string {
	if j.quality == "" {
		return "main"
	}
	return j.quality
}

func (j *Job) Status() valueobjects.JobStatus {
	return j.status
}

func (j *Job) Progress() float64 {
	return j.progress
}

func (j *Job) Error() string {
	return j.error
}

func (j *Job) Metadata() *valueobjects.VideoMetadata {
	return j.metadata
}

func (j *Job) SetMetadata(metadata *valueobjects.VideoMetadata) {
	j.metadata = metadata
}

func (j *Job) StartedAt() *time.Time {
	return j.startedAt
}

func (j *Job) CompletedAt() *time.Time {
	return j.completedAt
}

func (j *Job) CreatedAt() time.Time {
	return j.createdAt
}

func (j *Job) UpdatedAt() time.Time {
	return j.updatedAt
}

func (j *Job) Start() {
	now := time.Now().UTC()
	j.status = valueobjects.JobStatusRunning
	j.startedAt = &now
	j.updatedAt = now
}

func (j *Job) UpdateProgress(progress float64) {
	j.progress = progress
	j.updatedAt = time.Now().UTC()
}

func (j *Job) Complete(metadata *valueobjects.VideoMetadata) {
	now := time.Now().UTC()
	j.status = valueobjects.JobStatusCompleted
	j.progress = 100.0
	j.metadata = metadata
	j.completedAt = &now
	j.updatedAt = now
}

func (j *Job) Fail(errorMessage string) {
	now := time.Now().UTC()
	j.status = valueobjects.JobStatusFailed
	j.error = errorMessage
	j.completedAt = &now
	j.updatedAt = now
}

func (j *Job) IsCompleted() bool {
	return j.status == valueobjects.JobStatusCompleted
}

func (j *Job) IsFailed() bool {
	return j.status == valueobjects.JobStatusFailed
}

func (j *Job) IsRunning() bool {
	return j.status == valueobjects.JobStatusRunning
}

func (j *Job) IsPending() bool {
	return j.status == valueobjects.JobStatusPending
}

func (j *Job) Validate() error {
	if j.AssetID().Value() == "" {
		return fmt.Errorf("asset ID is required")
	}

	if j.VideoID().Value() == "" {
		return fmt.Errorf("video ID is required")
	}

	if j.Input() == "" {
		return fmt.Errorf("input is required")
	}

	if j.Type().IsTranscode() && j.Output() == "" {
		return fmt.Errorf("output is required for transcode jobs")
	}

	if j.Type().IsTranscode() && j.Format().String() == "" {
		return fmt.Errorf("format is required for transcode jobs")
	}

	return nil
}

func (j *Job) CreateCompletionEvent(success bool, metadata interface{}, errorMessage string) interface{} {
	return map[string]interface{}{
		"jobId":        j.ID().Value(),
		"assetId":      j.AssetID().Value(),
		"videoId":      j.VideoID().Value(),
		"jobType":      j.Type().String(),
		"format":       j.Format().String(),
		"success":      success,
		"metadata":     metadata,
		"errorMessage": errorMessage,
	}
}
