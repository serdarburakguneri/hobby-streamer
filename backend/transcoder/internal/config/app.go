package config

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AppConfig struct {
	SQS *SQSConfig
}

func NewAppConfig(ctx context.Context, configManager *config.Manager, log *logger.Logger) (*AppConfig, error) {
	sqsConfig, err := NewSQSConfig(ctx, configManager, log)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		SQS: sqsConfig,
	}, nil
}
