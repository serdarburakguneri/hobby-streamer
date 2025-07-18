package config

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type SQSConfig struct {
	Producer *SQSProducerConfig
	Consumer *SQSConsumerConfig
}

func NewSQSConfig(ctx context.Context, configManager *config.Manager, log *logger.Logger) (*SQSConfig, error) {
	producerConfig, err := NewSQSProducerConfig(ctx, configManager, log)
	if err != nil {
		return nil, err
	}

	consumerConfig, err := NewSQSConsumerConfig(ctx, configManager, producerConfig, log)
	if err != nil {
		return nil, err
	}

	return &SQSConfig{
		Producer: producerConfig,
		Consumer: consumerConfig,
	}, nil
}
