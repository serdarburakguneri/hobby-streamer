package sqs

import (
	"context"
	"encoding/json"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type MessageHandler func(ctx context.Context, msgType string, payload map[string]interface{}) error

type ConsumerRegistry struct {
	consumers []*ConsumerRegistryItem
	logger    *logger.Logger
}

type ConsumerRegistryItem struct {
	QueueURL string
	Handler  MessageHandler
	Consumer *Consumer
	Logger   *logger.Logger
}

func NewConsumerRegistry() *ConsumerRegistry {
	return &ConsumerRegistry{
		consumers: make([]*ConsumerRegistryItem, 0),
		logger:    logger.Get().WithService("consumer-registry"),
	}
}

func (r *ConsumerRegistry) Register(queueURL string, handler MessageHandler) {
	consumer := &ConsumerRegistryItem{
		QueueURL: queueURL,
		Handler:  handler,
		Logger:   r.logger.WithService("consumer"),
	}
	r.consumers = append(r.consumers, consumer)
}

func (r *ConsumerRegistry) Start(ctx context.Context) error {
	r.logger.Info("Starting consumer registry", "consumer_count", len(r.consumers))

	for _, consumer := range r.consumers {
		if consumer.QueueURL == "" {
			r.logger.Warn("Skipping consumer with empty queue URL")
			continue
		}

		sqsConsumer, err := NewConsumer(ctx, consumer.QueueURL)
		if err != nil {
			r.logger.WithError(err).Error("Failed to create SQS consumer", "queue_url", consumer.QueueURL)
			continue
		}

		consumer.Consumer = sqsConsumer
		r.logger.Info("Consumer initialized successfully", "queue_url", consumer.QueueURL)

		go func(c *ConsumerRegistryItem) {
			c.Consumer.Start(ctx, func(msg Message) error {
				var payload map[string]interface{}
				if err := json.Unmarshal(msg.Payload, &payload); err != nil {
					c.Logger.WithError(err).Error("Failed to unmarshal message payload")
					return err
				}
				c.Logger.Info("Processing message", "type", msg.Type, "payload", payload)
				return c.Handler(ctx, msg.Type, payload)
			})
		}(consumer)
	}

	return nil
}

func (r *ConsumerRegistry) Stop() {
	r.logger.Info("Stopping consumer registry")
}
