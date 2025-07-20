package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type EventConsumer struct {
	consumerRegistry *sqs.ConsumerRegistry
	logger           *logger.Logger
	assetService     AssetService
}

type AssetService interface {
	UpdateVideoAnalysis(ctx context.Context, assetID, videoID string, metadata *messages.JobCompletionPayload) error
	UpdateVideoTranscoding(ctx context.Context, assetID, videoID, format string, metadata *messages.JobCompletionPayload) error
}

func NewEventConsumer(assetService AssetService) *EventConsumer {
	return &EventConsumer{
		consumerRegistry: sqs.NewConsumerRegistry(),
		logger:           logger.WithService("sqs-event-consumer"),
		assetService:     assetService,
	}
}

func (c *EventConsumer) RegisterCompletionQueue(queueURL string) {
	c.consumerRegistry.Register(queueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		c.logger.Info("Job completion message received", "msgType", msgType, "payload", payload)

		if msgType != messages.MessageTypeJobCompleted {
			c.logger.Warn("Unknown message type in completion queue", "msgType", msgType)
			return nil
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		var jobCompletionPayload messages.JobCompletionPayload
		if err := json.Unmarshal(payloadBytes, &jobCompletionPayload); err != nil {
			return fmt.Errorf("failed to unmarshal job completion payload: %w", err)
		}

		switch jobCompletionPayload.JobType {
		case "analyze":
			return c.assetService.UpdateVideoAnalysis(ctx, jobCompletionPayload.AssetID, jobCompletionPayload.VideoID, &jobCompletionPayload)
		case "transcode":
			return c.assetService.UpdateVideoTranscoding(ctx, jobCompletionPayload.AssetID, jobCompletionPayload.VideoID, jobCompletionPayload.Format, &jobCompletionPayload)
		default:
			c.logger.Warn("Unknown job type in completion payload", "job_type", jobCompletionPayload.JobType)
			return nil
		}
	})
}

func (c *EventConsumer) Start(ctx context.Context) error {
	c.logger.Info("Starting SQS event consumer")
	return c.consumerRegistry.Start(ctx)
}

func (c *EventConsumer) Stop() error {
	c.logger.Info("Stopping SQS event consumer")
	c.consumerRegistry.Stop()
	return nil
}
