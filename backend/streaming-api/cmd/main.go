package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	sbootstrap "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/bootstrap"
	streamevents "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/events"
	httphandler "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/http"
)

func main() {
	configManager, secretsManager, cfg, dynamicCfg, err := sbootstrap.InitConfig("streaming-api")
	if err != nil {
		logger.Get().WithError(err).Error("Failed to initialize config")
		os.Exit(1)
	}
	defer configManager.Close()
	sbootstrap.InitLogger(cfg)
	log := logger.WithService(cfg.Service)
	log.Info("Starting streaming-api service", "environment", cfg.Environment)

	assetService, bucketService := sbootstrap.InitServices(cfg, dynamicCfg, secretsManager)

	handler := httphandler.NewHandler(assetService, bucketService, cfg)
	router := handler.SetupRoutes()
	wrapped := sbootstrap.InitRouter(router, cfg)
	server := sbootstrap.InitServer(wrapped, cfg)

	go func() {
		log.Info("Starting HTTP server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Failed to start server")
			os.Exit(1)
		}
	}()

	var stopInvalidator func()
	if cfg.Features.EnableCaching {
		cacheSvc, closeCache := sbootstrap.InitCacheService(cfg, dynamicCfg)
		if cacheSvc != nil {
			kafkaBootstrap := dynamicCfg.GetStringFromComponent("kafka", "bootstrap_servers")
			invalidator := streamevents.NewCacheInvalidator(cacheSvc)
			stopInvalidator = sbootstrap.InitCacheInvalidator(invalidator, kafkaBootstrap)
		}
		defer func() {
			if stopInvalidator != nil {
				stopInvalidator()
			}
			_ = closeCache()
		}()
	}

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
