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
		switch jobType {
		case "analyze":
			if videoMetadata, ok := metadata.(*domainjob.VideoMetadata); ok {
				payload.Width = videoMetadata.Width
				payload.Height = videoMetadata.Height
				payload.Duration = videoMetadata.Duration
				payload.Bitrate = videoMetadata.Bitrate
				payload.Codec = videoMetadata.Codec
				payload.Size = videoMetadata.Size
				payload.ContentType = videoMetadata.ContentType
			}
		case "transcode":
			if transcodeMetadata, ok := metadata.(*domainjob.TranscodeMetadata); ok {
				payload.Bucket = transcodeMetadata.Bucket
				payload.Key = transcodeMetadata.Key
				payload.URL = transcodeMetadata.OutputURL
				payload.Size = transcodeMetadata.Size
				payload.ContentType = transcodeMetadata.ContentType
				payload.Format = transcodeMetadata.Format
				payload.SegmentCount = transcodeMetadata.SegmentCount
				payload.VideoCodec = transcodeMetadata.VideoCodec
				payload.AudioCodec = transcodeMetadata.AudioCodec
				payload.AvgSegmentDuration = transcodeMetadata.AvgSegmentDuration
				payload.Segments = transcodeMetadata.Segments
			}
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
