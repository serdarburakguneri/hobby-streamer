package config

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AppConfig struct {
	Cache   *CacheConfig
	Service *ServiceConfig
	HTTP    *HTTPConfig
}

func NewAppConfig(configManager *config.Manager, secretsManager *config.SecretsManager, log *logger.Logger) (*AppConfig, error) {
	cacheConfig, err := NewCacheConfig(configManager, log)
	if err != nil {
		return nil, err
	}

	serviceConfig := NewServiceConfig(cacheConfig.Service, configManager, secretsManager)
	httpConfig := NewHTTPConfig(serviceConfig.StreamingService)

	return &AppConfig{
		Cache:   cacheConfig,
		Service: serviceConfig,
		HTTP:    httpConfig,
	}, nil
}

func (ac *AppConfig) Close() {
	if ac.Cache != nil {
		ac.Cache.Close()
	}
}
