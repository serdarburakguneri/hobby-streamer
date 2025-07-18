package config

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/cache"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/service"
)

type ServiceConfig struct {
	StreamingService *service.Service
}

func NewServiceConfig(cacheService *cache.Service, configManager *config.Manager, secretsManager *config.SecretsManager) *ServiceConfig {
	dynamicCfg := configManager.GetDynamicConfig()
	cfg := configManager.GetConfig()

	assetManagerURL := dynamicCfg.GetStringFromComponent("asset_manager", "url")
	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	realm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	clientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")
	clientSecret := secretsManager.Get("keycloak_client_secret")

	circuitBreakerConfig := errors.CircuitBreakerConfig{
		Threshold: int64(cfg.CircuitBreaker.Threshold),
		Timeout:   cfg.CircuitBreaker.Timeout,
	}

	streamingService := service.NewService(cacheService, assetManagerURL, keycloakURL, realm, clientID, clientSecret, circuitBreakerConfig)

	return &ServiceConfig{
		StreamingService: streamingService,
	}
}
