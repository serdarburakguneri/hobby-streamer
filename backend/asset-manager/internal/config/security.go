package config

import (
	"net/http"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
)

type SecurityConfig struct {
	Middleware func(http.Handler) http.Handler
}

func NewSecurityConfig(configManager *config.Manager, log *logger.Logger) *SecurityConfig {
	cfg := configManager.GetConfig()

	securityMiddleware := func(next http.Handler) http.Handler {
		handler := next

		handler = security.SecurityHeadersMiddleware()(handler)
		handler = security.RateLimitMiddleware(cfg.Security.RateLimit.Requests, cfg.Security.RateLimit.Window)(handler)
		handler = security.CORSMiddleware(
			cfg.Security.CORS.AllowedOrigins,
			cfg.Security.CORS.AllowedMethods,
			cfg.Security.CORS.AllowedHeaders,
		)(handler)
		handler = security.InputValidationMiddleware()(handler)
		handler = security.LoggingMiddleware(log)(handler)

		return handler
	}

	return &SecurityConfig{
		Middleware: securityMiddleware,
	}
}
