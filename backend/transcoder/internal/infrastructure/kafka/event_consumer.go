package kafka

import (
	"context"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	appjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/application/job"
)

type TranscoderEventConsumer struct {
	jobService appjob.JobApplicationService
	producer   *events.Producer
	consumer   *events.Consumer
	logger     *logger.Logger
}

func NewTranscoderEventConsumer(jobService appjob.JobApplicationService, producer *events.Producer) *TranscoderEventConsumer {
	return &TranscoderEventConsumer{
		jobService: jobService,
		producer:   producer,
		logger:     logger.WithService("transcoder-event-consumer"),
	}
}

func (c *TranscoderEventConsumer) Start(ctx context.Context, bootstrapServers string) error {
	cfg := events.DefaultConsumerConfig()
	cfg.BootstrapServers = []string{bootstrapServers}
	cfg.GroupID = events.TranscoderGroupID
	cfg.Topics = []string{events.AnalyzeJobRequestedTopic, events.HLSJobRequestedTopic, events.DASHJobRequestedTopic}

	consumer, err := events.NewConsumer(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	c.consumer = consumer

	consumer.Subscribe(events.AnalyzeJobRequestedTopic, c.HandleAnalyzeJobRequested)
	consumer.Subscribe(events.HLSJobRequestedTopic, c.HandleHLSJobRequested)
	consumer.Subscribe(events.DASHJobRequestedTopic, c.HandleDASHJobRequested)

	c.logger.Info("Starting Transcoder Kafka event consumer", "group_id", events.TranscoderGroupID, "topics", []string{events.AnalyzeJobRequestedTopic, events.HLSJobRequestedTopic})

	go func() {
		if err := consumer.Start(ctx); err != nil {
			c.logger.WithError(err).Error("Kafka consumer error")
		}
	}()

	return nil
}

func (c *TranscoderEventConsumer) Stop() error {
	if c.consumer != nil {
		c.logger.Info("Stopping Transcoder Kafka event consumer")
		return c.consumer.Stop()
	}
	return nil
}
