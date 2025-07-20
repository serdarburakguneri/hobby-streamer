package sqs

import (
	"context"
	"encoding/json"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	appjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/application/job"
)

type Consumer struct {
	consumerRegistry *sqs.ConsumerRegistry
	appService       appjob.JobApplicationService
	logger           *logger.Logger
}

func NewConsumer(appService appjob.JobApplicationService) *Consumer {
	return &Consumer{
		consumerRegistry: sqs.NewConsumerRegistry(),
		appService:       appService,
		logger:           logger.WithService("sqs-consumer"),
	}
}

func (c *Consumer) RegisterQueues(configManager *config.Manager) error {
	dynamicCfg := configManager.GetDynamicConfig()

	jobQueueURL := dynamicCfg.GetStringFromComponent("sqs", "job_queue_url")

	c.logger.Info("Registering SQS job queue", "job_queue", jobQueueURL)

	c.consumerRegistry.Register(jobQueueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		c.logger.Info("Job queue message received", "msgType", msgType, "payload", payload)

		if msgType != messages.MessageTypeJob {
			c.logger.Warn("Unknown message type in job queue", "msgType", msgType)
			return nil
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			c.logger.WithError(err).Error("Failed to marshal job payload")
			return errors.NewInternalError("failed to marshal job payload", err)
		}

		var jobPayload messages.JobPayload
		if err := json.Unmarshal(payloadBytes, &jobPayload); err != nil {
			c.logger.WithError(err).Error("Failed to unmarshal job payload")
			return errors.NewInternalError("failed to unmarshal job payload", err)
		}

		return c.appService.ProcessJob(ctx, jobPayload)
	})

	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("Starting SQS consumer")
	return c.consumerRegistry.Start(ctx)
}

func (c *Consumer) Stop() {
	c.logger.Info("Stopping SQS consumer")
	c.consumerRegistry.Stop()
}
