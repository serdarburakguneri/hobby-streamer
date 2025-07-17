package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/cache"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/handler"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/service"
)

func main() {
	configManager, err := config.NewManager("streaming-api")
	if err != nil {
		logger.Get().WithError(err).Error("Failed to initialize config")
		os.Exit(1)
	}
	defer configManager.Close()

	secretsManager := config.NewSecretsManager()
	secretsManager.LoadFromEnvironment()

	cfg := configManager.GetConfig()
	dynamicCfg := configManager.GetDynamicConfig()

	logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	log := logger.WithService(cfg.Service)
	log.Info("Starting streaming-api service", "environment", cfg.Environment)

	host := dynamicCfg.GetStringFromComponent("redis", "host")
	port := dynamicCfg.GetIntFromComponent("redis", "port")
	db := dynamicCfg.GetIntFromComponent("redis", "db")
	password := dynamicCfg.GetStringFromComponent("redis", "password")

	log.Info("Redis config", "host", host, "port", port, "db", db, "password_set", password != "")

	redisClient, err := cache.NewRedisClientWithConfig(host, port, db, password)
	if err != nil {
		log.WithError(err).Error("Failed to connect to Redis")
		os.Exit(1)
	}
	defer redisClient.Close()

	log.Info("Redis connection established", "host", host, "port", port, "db", db)

	cacheService := cache.NewService(redisClient)

	assetManagerURL := dynamicCfg.GetStringFromComponent("asset_manager", "url")
	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	realm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	clientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")
	clientSecret := secretsManager.Get("keycloak_client_secret")

	streamingService := service.NewService(cacheService, assetManagerURL, keycloakURL, realm, clientID, clientSecret)
	handler := handler.NewHandler(streamingService)

	router := handler.SetupRoutes()

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("Starting HTTP server", "port", cfg.Server.Port)
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
