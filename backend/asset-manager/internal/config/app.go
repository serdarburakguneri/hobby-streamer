package config

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AppConfig struct {
	Database *DatabaseConfig
	SQS      *SQSConfig
	Auth     *AuthConfig
	CORS     *CORSConfig
	Logging  *LoggingConfig
	GraphQL  *GraphQLConfig
	Services *ServicesConfig
}

func NewAppConfig(ctx context.Context, configManager *config.Manager, secretsManager *config.SecretsManager, log *logger.Logger) (*AppConfig, error) {
	databaseConfig, err := NewDatabaseConfig(configManager, secretsManager, log)
	if err != nil {
		return nil, err
	}

	authConfig, err := NewAuthConfig(configManager)
	if err != nil {
		databaseConfig.Close()
		return nil, err
	}

	corsConfig := NewCORSConfig()
	loggingConfig := NewLoggingConfig(log)

	dynamicCfg := configManager.GetDynamicConfig()
	servicesConfig := NewServicesConfig(databaseConfig.Driver, nil, dynamicCfg)

	sqsConfig, err := NewSQSConfig(ctx, configManager, servicesConfig.AssetService, log)
	if err != nil {
		databaseConfig.Close()
		return nil, err
	}

	servicesConfig = NewServicesConfig(databaseConfig.Driver, sqsConfig.Producer, dynamicCfg)

	graphQLConfig := NewGraphQLConfig(servicesConfig.AssetService, servicesConfig.BucketService)

	return &AppConfig{
		Database: databaseConfig,
		SQS:      sqsConfig,
		Auth:     authConfig,
		CORS:     corsConfig,
		Logging:  loggingConfig,
		GraphQL:  graphQLConfig,
		Services: servicesConfig,
	}, nil
}

func (ac *AppConfig) Close() {
	if ac.Database != nil {
		ac.Database.Close()
	}
}
