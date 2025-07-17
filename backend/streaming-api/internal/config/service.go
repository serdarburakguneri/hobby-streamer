package config

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/cache"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/service"
)

type ServiceConfig struct {
	StreamingService *service.Service
}

func NewServiceConfig(cacheService *cache.Service, configManager *config.Manager, secretsManager *config.SecretsManager) *ServiceConfig {
	dynamicCfg := configManager.GetDynamicConfig()

	assetManagerURL := dynamicCfg.GetStringFromComponent("asset_manager", "url")
	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	realm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	clientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")
	clientSecret := secretsManager.Get("keycloak_client_secret")

	streamingService := service.NewService(cacheService, assetManagerURL, keycloakURL, realm, clientID, clientSecret)

	return &ServiceConfig{
		StreamingService: streamingService,
	}
}
