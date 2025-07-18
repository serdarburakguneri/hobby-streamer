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
	Security *SecurityConfig
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

	securityConfig := NewSecurityConfig(configManager, log)
	loggingConfig := NewLoggingConfig(log)

	dynamicCfg := configManager.GetDynamicConfig()
	servicesConfig := NewServicesConfig(databaseConfig.Driver, dynamicCfg)

	sqsConfig, err := NewSQSConfig(ctx, configManager, servicesConfig.AssetService, log)
	if err != nil {
		databaseConfig.Close()
		return nil, err
	}

	graphQLConfig := NewGraphQLConfig(servicesConfig.AssetService, servicesConfig.BucketService, configManager.GetConfig().Security.CORS.AllowedOrigins)

	return &AppConfig{
		Database: databaseConfig,
		SQS:      sqsConfig,
		Auth:     authConfig,
		Security: securityConfig,
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
