package kafka

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
)

type KafkaProducer struct {
	producer *events.Producer
}

func NewKafkaProducer(producer *events.Producer) *KafkaProducer {
	return &KafkaProducer{
		producer: producer,
	}
}

func (k *KafkaProducer) SendEvent(ctx context.Context, topic string, event *events.Event) error {
	return k.producer.SendEvent(ctx, topic, event)
}
