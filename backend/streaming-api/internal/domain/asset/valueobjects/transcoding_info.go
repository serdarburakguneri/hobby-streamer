package valueobjects

import "time"

type TranscodingInfo struct {
	JobID       *string
	Progress    *float64
	OutputURL   *string
	Error       *string
	CompletedAt *time.Time
}
