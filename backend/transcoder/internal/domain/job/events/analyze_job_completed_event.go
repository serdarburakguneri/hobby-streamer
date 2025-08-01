package events

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
)

type AnalyzeJobCompletedEvent struct {
	JobCompletedBase
	URL         string  `json:"url,omitempty"`
	Bucket      string  `json:"bucket,omitempty"`
	Key         string  `json:"key,omitempty"`
	Width       int     `json:"width,omitempty"`
	Height      int     `json:"height,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	Bitrate     int     `json:"bitrate,omitempty"`
	Codec       string  `json:"codec,omitempty"`
	Size        int64   `json:"size,omitempty"`
	ContentType string  `json:"contentType,omitempty"`
}

func (*AnalyzeJobCompletedEvent) Topic() string          { return events.AnalyzeJobCompletedTopic }
func (*AnalyzeJobCompletedEvent) CloudEventType() string { return events.JobAnalyzeCompletedEventType }
func (e *AnalyzeJobCompletedEvent) Type() string         { return "job.analyze.completed" }
func (e *AnalyzeJobCompletedEvent) Data() interface{}    { return e }
