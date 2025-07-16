package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/cache"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/handler"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/service"
)

func main() {
	log := logger.Get().WithService("streaming-api")
	log.Info("Starting streaming-api service")

	redisClient, err := cache.NewRedisClient()
	if err != nil {
		log.WithError(err).Error("Failed to connect to Redis")
		os.Exit(1)
	}
	defer redisClient.Close()

	cacheService := cache.NewService(redisClient)
	streamingService := service.NewService(cacheService)
	handler := handler.NewHandler(streamingService)

	router := handler.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Starting HTTP server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
		os.Exit(1)
	}

	log.Info("Server exited")
}
