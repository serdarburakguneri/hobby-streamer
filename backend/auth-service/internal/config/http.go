package config

import (
	"net/http"

	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/auth"
	httphandler "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/http"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type HTTPConfig struct {
	Handler http.Handler
}

func NewHTTPConfig(authService *auth.Service, log *logger.Logger) *HTTPConfig {
	router := httphandler.NewRouter(authService)
	handler := logger.RequestLoggingMiddleware(log)(httphandler.CORS(router))

	return &HTTPConfig{
		Handler: handler,
	}
}
