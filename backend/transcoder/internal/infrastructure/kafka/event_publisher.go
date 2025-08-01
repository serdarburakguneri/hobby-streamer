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

func (p *KafkaEventPublisher) PublishJobCompleted(ctx context.Context, ev jobevents.CompletedEvent) error {
	ce := events.NewEvent(ev.CloudEventType(), ev.Data()).
		SetSource("transcoder").
		AddExtension("subject", ev.ID())
	if err := p.producer.SendEvent(ctx, ev.Topic(), ce); err != nil {
		p.logger.WithError(err).Error("Failed to publish job completion event", "topic", ev.Topic(), "job_id", ev.ID())
		return err
	}
	p.logger.Info("Published job completion event", "topic", ev.Topic(), "job_id", ev.ID())
	return nil
}
