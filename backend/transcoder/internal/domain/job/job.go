package job

import (
	"time"
)

type JobType string
type JobStatus string
type JobFormat string

const (
	JobTypeAnalyze   JobType = "analyze"
	JobTypeTranscode JobType = "transcode"

	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"

	JobFormatHLS  JobFormat = "hls"
	JobFormatDASH JobFormat = "dash"
)

type Job struct {
	id          string
	jobType     JobType
	format      JobFormat
	assetID     string
	videoID     string
	input       string
	output      string
	status      JobStatus
	progress    float64
	error       string
	metadata    *VideoMetadata
	startedAt   *time.Time
	completedAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewAnalyzeJob(assetID, videoID, input string) *Job {
	now := time.Now().UTC()
	return &Job{
		id:        generateJobID(),
		jobType:   JobTypeAnalyze,
		assetID:   assetID,
		videoID:   videoID,
		input:     input,
		status:    JobStatusPending,
		progress:  0.0,
		createdAt: now,
		updatedAt: now,
	}
}

func NewTranscodeJob(assetID, videoID, input, output string, format JobFormat) *Job {
	now := time.Now().UTC()
	return &Job{
		id:        generateJobID(),
		jobType:   JobTypeTranscode,
		format:    format,
		assetID:   assetID,
		videoID:   videoID,
		input:     input,
		output:    output,
		status:    JobStatusPending,
		progress:  0.0,
		createdAt: now,
		updatedAt: now,
	}
}

func (j *Job) ID() string {
	return j.id
}

func (j *Job) Type() JobType {
	return j.jobType
}

func (j *Job) Format() JobFormat {
	return j.format
}

func (j *Job) AssetID() string {
	return j.assetID
}

func (j *Job) VideoID() string {
	return j.videoID
}

func (j *Job) Input() string {
	return j.input
}

func (j *Job) Output() string {
	return j.output
}

func (j *Job) Status() JobStatus {
	return j.status
}

func (j *Job) Progress() float64 {
	return j.progress
}

func (j *Job) Error() string {
	return j.error
}

func (j *Job) Metadata() *VideoMetadata {
	return j.metadata
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
	j.status = JobStatusRunning
	j.startedAt = &now
	j.updatedAt = now
}

func (j *Job) UpdateProgress(progress float64) {
	j.progress = progress
	j.updatedAt = time.Now().UTC()
}

func (j *Job) Complete(metadata *VideoMetadata) {
	now := time.Now().UTC()
	j.status = JobStatusCompleted
	j.progress = 100.0
	j.metadata = metadata
	j.completedAt = &now
	j.updatedAt = now
}

func (j *Job) Fail(errorMessage string) {
	now := time.Now().UTC()
	j.status = JobStatusFailed
	j.error = errorMessage
	j.completedAt = &now
	j.updatedAt = now
}

func (j *Job) IsCompleted() bool {
	return j.status == JobStatusCompleted
}

func (j *Job) IsFailed() bool {
	return j.status == JobStatusFailed
}

func (j *Job) IsRunning() bool {
	return j.status == JobStatusRunning
}

func (j *Job) IsPending() bool {
	return j.status == JobStatusPending
}
