package events

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

func NewAnalyzeJobCompletedEvent(job *entity.Job, success bool, metadata interface{}, errorMessage string) CompletedEvent {
	ev := &AnalyzeJobCompletedEvent{
		JobCompletedBase: JobCompletedBase{
			JobID:        job.ID().Value(),
			AssetID:      job.AssetID().Value(),
			VideoID:      job.VideoID().Value(),
			Success:      success,
			ErrorMessage: errorMessage,
			CompletedAt:  time.Now().UTC().Format(time.RFC3339),
		},
	}
	if success && metadata != nil {
		if m, ok := metadata.(*valueobjects.TranscodeMetadata); ok {
			ev.URL = m.OutputURL
			ev.Bucket = m.Bucket
			ev.Key = m.Key
			ev.Width = m.Width
			ev.Height = m.Height
			ev.Duration = m.Duration
			ev.Bitrate = m.Bitrate
			ev.Codec = m.Codec
			ev.Size = m.Size
			ev.ContentType = m.ContentType
		}
	}
	return ev
}

func NewHLSJobCompletedEvent(job *entity.Job, success bool, metadata interface{}, errorMessage string) CompletedEvent {
	ev := &HLSJobCompletedEvent{
		JobCompletedBase: JobCompletedBase{
			JobID:        job.ID().Value(),
			AssetID:      job.AssetID().Value(),
			VideoID:      job.VideoID().Value(),
			Success:      success,
			ErrorMessage: errorMessage,
			CompletedAt:  time.Now().UTC().Format(time.RFC3339),
		},
		Format: "hls",
	}
	if success && metadata != nil {
		if m, ok := metadata.(*valueobjects.TranscodeMetadata); ok {
			ev.URL = m.OutputURL
			ev.Bucket = m.Bucket
			ev.Key = m.Key
			ev.Width = m.Width
			ev.Height = m.Height
			ev.Duration = m.Duration
			ev.Bitrate = m.Bitrate
			ev.ContentType = m.ContentType
			ev.SegmentCount = m.SegmentCount
			ev.VideoCodec = m.VideoCodec
			ev.AudioCodec = m.AudioCodec
			ev.AvgSegmentDuration = m.AvgSegmentDuration
			ev.Segments = m.Segments
			ev.FrameRate = m.FrameRate
			ev.AudioChannels = m.AudioChannels
			ev.AudioSampleRate = m.AudioSampleRate
		}
	}
	return ev
}

func NewDASHJobCompletedEvent(job *entity.Job, success bool, metadata interface{}, errorMessage string) CompletedEvent {
	ev := &DASHJobCompletedEvent{
		JobCompletedBase: JobCompletedBase{
			JobID:        job.ID().Value(),
			AssetID:      job.AssetID().Value(),
			VideoID:      job.VideoID().Value(),
			Success:      success,
			ErrorMessage: errorMessage,
			CompletedAt:  time.Now().UTC().Format(time.RFC3339),
		},
		Format: "dash",
	}
	if success && metadata != nil {
		if m, ok := metadata.(*valueobjects.TranscodeMetadata); ok {
			ev.URL = m.OutputURL
			ev.Bucket = m.Bucket
			ev.Key = m.Key
			ev.Width = m.Width
			ev.Height = m.Height
			ev.Duration = m.Duration
			ev.Bitrate = m.Bitrate
			ev.ContentType = m.ContentType
			ev.SegmentCount = m.SegmentCount
			ev.VideoCodec = m.VideoCodec
			ev.AudioCodec = m.AudioCodec
			ev.AvgSegmentDuration = m.AvgSegmentDuration
			ev.Segments = m.Segments
			ev.FrameRate = m.FrameRate
			ev.AudioChannels = m.AudioChannels
			ev.AudioSampleRate = m.AudioSampleRate
		}
	}
	return ev
}

var builderMap = map[string]func(*entity.Job, bool, interface{}, string) CompletedEvent{
	"analyze":        NewAnalyzeJobCompletedEvent,
	"transcode:hls":  NewHLSJobCompletedEvent,
	"transcode:dash": NewDASHJobCompletedEvent,
}

func BuildCompletedEvent(job *entity.Job, success bool, metadata interface{}, errorMessage string) CompletedEvent {
	key := job.Type().String()
	if !job.Type().IsAnalyze() {
		key += ":" + job.Format().String()
	}
	if builder, ok := builderMap[key]; ok {
		return builder(job, success, metadata, errorMessage)
	}
	return NewAnalyzeJobCompletedEvent(job, success, metadata, errorMessage)
}
