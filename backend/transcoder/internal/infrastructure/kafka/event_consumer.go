package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	appjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/application/job"
)

type JobAnalyzeRequestedEvent struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Input   string `json:"input"`
}

type HLSJobRequestedEvent struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Input   string `json:"input"`
}

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
	consumerConfig := &events.ConsumerConfig{
		BootstrapServers:  []string{bootstrapServers},
		GroupID:           events.TranscoderGroupID,
		Topics:            []string{events.AnalyzeJobRequestedTopic, events.HLSJobRequestedTopic},
		SessionTimeout:    10 * time.Second,
		HeartbeatInterval: 3 * time.Second,
	}

	consumer, err := events.NewConsumer(ctx, consumerConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	c.consumer = consumer

	consumer.Subscribe(events.AnalyzeJobRequestedTopic, c.HandleAnalyzeJobRequested)
	consumer.Subscribe(events.HLSJobRequestedTopic, c.HandleHLSJobRequested)

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

func (c *TranscoderEventConsumer) unmarshalEventData(event *events.Event, target interface{}) error {
	var dataBytes []byte
	switch v := event.Data.(type) {
	case []byte:
		dataBytes = v
	case map[string]interface{}:
		var err error
		dataBytes, err = json.Marshal(v)
		if err != nil {
			c.logger.WithError(err).Error("Failed to marshal event data")
			return err
		}
	default:
		c.logger.Error("Unexpected event data type", "type", fmt.Sprintf("%T", event.Data))
		return fmt.Errorf("unexpected event data type: %T", event.Data)
	}

	return json.Unmarshal(dataBytes, target)
}

func (c *TranscoderEventConsumer) HandleAnalyzeJobRequested(ctx context.Context, event *events.Event) error {
	c.logger.Info("Analyze job requested event received", "event_id", event.ID, "source", event.Source)

	var e JobAnalyzeRequestedEvent
	if err := c.unmarshalEventData(event, &e); err != nil {
		c.logger.WithError(err).Error("Failed to unmarshal analyze job event")
		return err
	}

	payload := messages.JobPayload{
		JobType: "analyze",
		AssetID: e.AssetID,
		VideoID: e.VideoID,
		Input:   e.Input,
	}

	if err := c.jobService.ProcessJob(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to process analyze job", "asset_id", e.AssetID, "video_id", e.VideoID)
		return err
	}

	c.logger.Info("Analyze job processed successfully", "asset_id", e.AssetID, "video_id", e.VideoID)
	return nil
}

func (c *TranscoderEventConsumer) HandleHLSJobRequested(ctx context.Context, event *events.Event) error {
	c.logger.Info("HLS job requested event received", "event_id", event.ID, "source", event.Source)

	var e HLSJobRequestedEvent
	if err := c.unmarshalEventData(event, &e); err != nil {
		c.logger.WithError(err).Error("Failed to unmarshal HLS job event")
		return err
	}

	payload := messages.JobPayload{
		JobType: "transcode",
		AssetID: e.AssetID,
		VideoID: e.VideoID,
		Input:   e.Input,
		Format:  "hls",
		Quality: "main",
	}

	if err := c.jobService.ProcessJob(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to process HLS job", "asset_id", e.AssetID, "video_id", e.VideoID)
		return err
	}

	c.logger.Info("HLS job processed successfully", "asset_id", e.AssetID, "video_id", e.VideoID)
	return nil
}
