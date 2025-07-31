package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type EventHandler func(ctx context.Context, event *Event) error

type Consumer struct {
	consumer sarama.ConsumerGroup
	logger   *logger.Logger
	handlers map[string]EventHandler
	topics   []string
	groupID  string
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

type ConsumerConfig struct {
	BootstrapServers  []string
	GroupID           string
	Topics            []string
	AutoOffsetReset   string
	SessionTimeout    time.Duration
	HeartbeatInterval time.Duration
}

func DefaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		AutoOffsetReset:   "earliest",
		SessionTimeout:    10 * time.Second,
		HeartbeatInterval: 3 * time.Second,
	}
}

func NewConsumer(ctx context.Context, config *ConsumerConfig) (*Consumer, error) {
	if config == nil {
		config = DefaultConsumerConfig()
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Group.Session.Timeout = config.SessionTimeout
	saramaConfig.Consumer.Group.Heartbeat.Interval = config.HeartbeatInterval
	saramaConfig.Version = sarama.V2_8_1_0

	consumer, err := sarama.NewConsumerGroup(config.BootstrapServers, config.GroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Consumer{
		consumer: consumer,
		logger:   logger.WithService("kafka-consumer"),
		handlers: make(map[string]EventHandler),
		topics:   config.Topics,
		groupID:  config.GroupID,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (c *Consumer) Subscribe(topic string, handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[topic] = handler
	c.logger.Info("Subscribed to topic", "topic", topic, "group_id", c.groupID)
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("Starting Kafka consumer", "group_id", c.groupID, "topics", c.topics)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := c.consumer.Consume(ctx, c.topics, c)
			if err != nil {
				c.logger.WithError(err).Error("Error from consumer", "group_id", c.groupID)
				time.Sleep(time.Second)
			}
		}
	}
}

func (c *Consumer) Stop() error {
	c.logger.Info("Stopping Kafka consumer", "group_id", c.groupID)
	c.cancel()
	return c.consumer.Close()
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session setup", "group_id", c.groupID)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session cleanup", "group_id", c.groupID)
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.logger.Info("Starting to consume claims", "topic", claim.Topic(), "partition", claim.Partition())

	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				continue
			}

			start := time.Now()
			err := c.processMessage(session.Context(), message)
			duration := time.Since(start)

			if err != nil {
				c.logger.WithError(err).Error("Failed to process message",
					"topic", message.Topic,
					"partition", message.Partition,
					"offset", message.Offset,
					"duration_ms", duration.Milliseconds(),
				)
				continue
			}

			session.MarkMessage(message, "")
			c.logger.Debug("Message processed successfully",
				"topic", message.Topic,
				"partition", message.Partition,
				"offset", message.Offset,
				"duration_ms", duration.Milliseconds(),
			)

		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var event Event
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	c.mu.RLock()
	handler, exists := c.handlers[message.Topic]
	c.mu.RUnlock()

	if !exists {
		c.logger.Warn("No handler registered for topic", "topic", message.Topic, "event_type", event.Type)
		return nil
	}

	c.logger.Debug("Processing event",
		"topic", message.Topic,
		"event_id", event.ID,
		"event_type", event.Type,
		"partition", message.Partition,
		"offset", message.Offset,
	)

	return handler(ctx, &event)
}

func (c *Consumer) GetConsumerGroup() sarama.ConsumerGroup {
	return c.consumer
}
