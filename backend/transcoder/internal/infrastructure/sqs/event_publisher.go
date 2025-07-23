package sqs

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	domainjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
)

type EventPublisher struct {
	completionProducer *sqs.Producer
	logger             *logger.Logger
}

func NewEventPublisher(completionProducer *sqs.Producer) *EventPublisher {
	return &EventPublisher{
		completionProducer: completionProducer,
		logger:             logger.WithService("sqs-event-publisher"),
	}
}

func (p *EventPublisher) PublishJobCompleted(ctx context.Context, jobType, assetID, videoID string, success bool, metadata interface{}, errorMessage string) error {
	payload := messages.JobCompletionPayload{
		JobType: jobType,
		AssetID: assetID,
		VideoID: videoID,
		Success: success,
		Error:   errorMessage,
	}

	if success && metadata != nil {
		if meta, ok := metadata.(*domainjob.TranscodeMetadata); ok {
			payload.Width = meta.Width
			payload.Height = meta.Height
			payload.Duration = meta.Duration
			payload.Bitrate = meta.Bitrate
			payload.Codec = meta.VideoCodec
			payload.Size = meta.Size
			payload.ContentType = meta.ContentType
			payload.Bucket = meta.Bucket
			payload.Key = meta.Key
			payload.URL = meta.OutputURL
			payload.Format = meta.Format
			payload.SegmentCount = meta.SegmentCount
			payload.VideoCodec = meta.VideoCodec
			payload.AudioCodec = meta.AudioCodec
			payload.AvgSegmentDuration = meta.AvgSegmentDuration
			payload.Segments = meta.Segments
			payload.FrameRate = meta.FrameRate
			payload.AudioChannels = meta.AudioChannels
			payload.AudioSampleRate = meta.AudioSampleRate
		}
	}

	err := p.completionProducer.SendMessage(ctx, messages.MessageTypeJobCompleted, payload)
	if err != nil {
		p.logger.WithError(err).Error("Failed to publish job completion event", "job_type", jobType, "asset_id", assetID, "video_id", videoID)
		return errors.NewExternalError("failed to publish job completion", err)
	}

	p.logger.Info("Job completion event published", "job_type", jobType, "asset_id", assetID, "video_id", videoID, "success", success)
	return nil
}
