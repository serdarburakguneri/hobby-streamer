package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *logger.Logger
	source   string
}

type ProducerConfig struct {
	BootstrapServers []string
	Source           string
	Compression      string
	RequiredAcks     sarama.RequiredAcks
	MaxMessageBytes  int
}

func DefaultProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		Compression:     "snappy",
		RequiredAcks:    sarama.WaitForAll,
		MaxMessageBytes: 1000000,
	}
}

func NewProducer(ctx context.Context, config *ProducerConfig) (*Producer, error) {
	if config == nil {
		config = DefaultProducerConfig()
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = config.RequiredAcks
	saramaConfig.Producer.MaxMessageBytes = config.MaxMessageBytes

	switch config.Compression {
	case "snappy":
		saramaConfig.Producer.Compression = sarama.CompressionSnappy
	case "gzip":
		saramaConfig.Producer.Compression = sarama.CompressionGZIP
	case "lz4":
		saramaConfig.Producer.Compression = sarama.CompressionLZ4
	default:
		saramaConfig.Producer.Compression = sarama.CompressionNone
	}

	saramaConfig.Version = sarama.V2_8_1_0

	producer, err := sarama.NewSyncProducer(config.BootstrapServers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{
		producer: producer,
		logger:   logger.WithService("kafka-producer"),
		source:   config.Source,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, event *Event) error {
	return p.SendEvent(ctx, topic, event)
}

func (p *Producer) SendEvent(ctx context.Context, topic string, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.Source == "" {
		event.SetSource(p.source)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(eventBytes),
		Key:   sarama.StringEncoder(p.getPartitionKey(event)),
		Headers: []sarama.RecordHeader{
			{Key: []byte("content-type"), Value: []byte("application/cloudevents+json")},
			{Key: []byte("event-id"), Value: []byte(event.ID)},
			{Key: []byte("event-type"), Value: []byte(event.Type)},
		},
	}

	if event.CorrelationID != "" {
		message.Headers = append(message.Headers, sarama.RecordHeader{
			Key: []byte("correlation-id"), Value: []byte(event.CorrelationID),
		})
	}

	start := time.Now()
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		p.logger.WithError(err).Error("Failed to send event to Kafka",
			"topic", topic,
			"event_id", event.ID,
			"event_type", event.Type,
			"partition", partition,
			"offset", offset,
		)
		return fmt.Errorf("failed to send event to Kafka: %w", err)
	}

	duration := time.Since(start)
	p.logger.Info("Event sent successfully",
		"topic", topic,
		"event_id", event.ID,
		"event_type", event.Type,
		"partition", partition,
		"offset", offset,
		"duration_ms", duration.Milliseconds(),
		"message_size_bytes", len(eventBytes),
	)

	return nil
}

func (p *Producer) SendEventWithKey(ctx context.Context, topic string, event *Event, partitionKey string) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.Source == "" {
		event.SetSource(p.source)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(eventBytes),
		Key:   sarama.StringEncoder(partitionKey),
		Headers: []sarama.RecordHeader{
			{Key: []byte("content-type"), Value: []byte("application/cloudevents+json")},
			{Key: []byte("event-id"), Value: []byte(event.ID)},
			{Key: []byte("event-type"), Value: []byte(event.Type)},
		},
	}

	if event.CorrelationID != "" {
		message.Headers = append(message.Headers, sarama.RecordHeader{
			Key: []byte("correlation-id"), Value: []byte(event.CorrelationID),
		})
	}

	start := time.Now()
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		p.logger.WithError(err).Error("Failed to send event to Kafka",
			"topic", topic,
			"event_id", event.ID,
			"event_type", event.Type,
			"partition_key", partitionKey,
			"partition", partition,
			"offset", offset,
		)
		return fmt.Errorf("failed to send event to Kafka: %w", err)
	}

	duration := time.Since(start)
	p.logger.Info("Event sent successfully",
		"topic", topic,
		"event_id", event.ID,
		"event_type", event.Type,
		"partition_key", partitionKey,
		"partition", partition,
		"offset", offset,
		"duration_ms", duration.Milliseconds(),
		"message_size_bytes", len(eventBytes),
	)

	return nil
}

func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		p.logger.WithError(err).Error("Failed to close Kafka producer")
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}
	p.logger.Info("Kafka producer closed successfully")
	return nil
}

func (p *Producer) getPartitionKey(event *Event) string {
	if event.Data != nil {
		if dataMap, ok := event.Data.(map[string]interface{}); ok {
			if assetID, exists := dataMap["assetId"]; exists {
				if assetIDStr, ok := assetID.(string); ok {
					return assetIDStr
				}
			}
			if bucketID, exists := dataMap["bucketId"]; exists {
				if bucketIDStr, ok := bucketID.(string); ok {
					return bucketIDStr
				}
			}
		}
	}

	return event.ID
}
