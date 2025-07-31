package kafka

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	jobevents "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/events"
)

type KafkaEventPublisher struct {
	producer *events.Producer
	logger   *logger.Logger
}

func NewKafkaEventPublisher(producer *events.Producer) *KafkaEventPublisher {
	return &KafkaEventPublisher{
		producer: producer,
		logger:   logger.WithService("kafka-event-publisher"),
	}
}

func (p *KafkaEventPublisher) PublishJobCompleted(ctx context.Context, event *jobevents.JobCompletedEvent) error {
	var topic string
	switch event.JobType {
	case "analyze":
		topic = events.AnalyzeJobCompletedTopic
	default:
		topic = events.DASHJobCompletedTopic
	}
	var ceType string
	switch event.JobType {
	case "analyze":
		ceType = events.JobAnalyzeCompletedEventType
	default:
		ceType = events.JobTranscodeCompletedEventType
	}
	ce := events.NewEvent(ceType, event).
		SetSource("transcoder").
		AddExtension("subject", event.JobID)
	if err := p.producer.SendEvent(ctx, topic, ce); err != nil {
		p.logger.WithError(err).Error("Failed to publish job completion event", "topic", topic, "job_id", event.JobID)
		return err
	}
	p.logger.Info("Published job completion event", "topic", topic, "job_id", event.JobID, "success", event.Success)
	return nil
}
