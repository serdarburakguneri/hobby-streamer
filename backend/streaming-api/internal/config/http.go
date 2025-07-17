package config

import (
	"net/http"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/handler"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/service"
)

type HTTPConfig struct {
	Handler http.Handler
}

func NewHTTPConfig(streamingService *service.Service) *HTTPConfig {
	handler := handler.NewHandler(streamingService)
	router := handler.SetupRoutes()

	return &HTTPConfig{
		Handler: router,
	}
}
