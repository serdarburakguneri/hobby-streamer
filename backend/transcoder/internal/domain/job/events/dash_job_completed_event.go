package events

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
)

type DASHJobCompletedEvent struct {
	JobCompletedBase
	Format             string   `json:"format"`
	URL                string   `json:"url,omitempty"`
	Bucket             string   `json:"bucket,omitempty"`
	Key                string   `json:"key,omitempty"`
	Width              int      `json:"width,omitempty"`
	Height             int      `json:"height,omitempty"`
	Duration           float64  `json:"duration,omitempty"`
	Bitrate            int      `json:"bitrate,omitempty"`
	ContentType        string   `json:"contentType,omitempty"`
	SegmentCount       int      `json:"segmentCount,omitempty"`
	VideoCodec         string   `json:"videoCodec,omitempty"`
	AudioCodec         string   `json:"audioCodec,omitempty"`
	AvgSegmentDuration float64  `json:"avgSegmentDuration,omitempty"`
	Segments           []string `json:"segments,omitempty"`
	FrameRate          string   `json:"frameRate,omitempty"`
	AudioChannels      int      `json:"audioChannels,omitempty"`
	AudioSampleRate    int      `json:"audioSampleRate,omitempty"`
}

func (*DASHJobCompletedEvent) Topic() string          { return events.DASHJobCompletedTopic }
func (*DASHJobCompletedEvent) CloudEventType() string { return events.JobTranscodeCompletedEventType }
func (e *DASHJobCompletedEvent) Type() string         { return "job.transcode.completed" }
func (e *DASHJobCompletedEvent) Data() interface{}    { return e }
