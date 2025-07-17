package config

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
)

type AuthConfig struct {
	Service *auth.Service
}

func NewAuthConfig(configManager *config.Manager, secretsManager *config.SecretsManager) *AuthConfig {
	dynamicCfg := configManager.GetDynamicConfig()

	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	realm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	clientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")
	clientSecret := dynamicCfg.GetStringFromComponent("keycloak", "client_secret")

	authService := auth.NewService(keycloakURL, realm, clientID, clientSecret)

	return &AuthConfig{
		Service: authService,
	}
}
