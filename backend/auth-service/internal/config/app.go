package config

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AppConfig struct {
	Auth     *AuthConfig
	HTTP     *HTTPConfig
	Security *SecurityConfig
}

func NewAppConfig(configManager *config.Manager, secretsManager *config.SecretsManager, log *logger.Logger) *AppConfig {
	authConfig := NewAuthConfig(configManager, secretsManager)
	httpConfig := NewHTTPConfig(authConfig.Service, log)
	securityConfig := NewSecurityConfig(configManager, log)

	return &AppConfig{
		Auth:     authConfig,
		HTTP:     httpConfig,
		Security: securityConfig,
	}
}
