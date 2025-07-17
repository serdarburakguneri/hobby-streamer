package config

import (
	"net/http"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type LoggingConfig struct {
	Middleware func(http.Handler) http.Handler
}

func NewLoggingConfig(log *logger.Logger) *LoggingConfig {
	loggerMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	return &LoggingConfig{
		Middleware: loggerMiddleware,
	}
}
