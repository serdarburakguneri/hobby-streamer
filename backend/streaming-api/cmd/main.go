package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	assetapp "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/application/asset"
	bucketapp "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/application/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/cache"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/graphql"
	httphandler "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/http"
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

	if cfg.Log.Async.Enabled {
		logger.InitAsync(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format, cfg.Log.Async.BufferSize)
	} else {
		logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	}
	log := logger.WithService(cfg.Service)
	log.Info("Starting streaming-api service", "environment", cfg.Environment)

	dynamicCfg := configManager.GetDynamicConfig()
	assetManagerURL := dynamicCfg.GetStringFromComponent("asset_manager", "url")
	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	realm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	clientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")
	clientSecret := secretsManager.Get("keycloak_client_secret")

	redisHost := dynamicCfg.GetStringFromComponent("redis", "host")
	redisPort := dynamicCfg.GetIntFromComponent("redis", "port")
	redisDB := dynamicCfg.GetIntFromComponent("redis", "db")
	redisPassword := dynamicCfg.GetStringFromComponent("redis", "password")

	circuitBreaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
		Name:      "asset-manager",
		Threshold: int64(cfg.CircuitBreaker.Threshold),
		Timeout:   cfg.CircuitBreaker.Timeout,
		OnStateChange: func(name string, from, to errors.CircuitState) {
			log.Info("Circuit breaker state changed", "name", name, "from", from, "to", to)
		},
	})

	redisClient, err := cache.NewRedisClientWithConfig(redisHost, redisPort, redisDB, redisPassword)
	if err != nil {
		log.WithError(err).Error("Failed to initialize Redis client")
		os.Exit(1)
	}
	defer redisClient.Close()

	cacheService := cache.NewService(redisClient, cache.TTLConfig{
		Bucket:      dynamicCfg.GetDurationFromComponent("redis", "ttl.bucket", 30*time.Minute),
		BucketsList: dynamicCfg.GetDurationFromComponent("redis", "ttl.buckets_list", 15*time.Minute),
		Asset:       dynamicCfg.GetDurationFromComponent("redis", "ttl.asset", 30*time.Minute),
		AssetsList:  dynamicCfg.GetDurationFromComponent("redis", "ttl.assets_list", 15*time.Minute),
	})

	serviceClient := auth.NewServiceClient(keycloakURL, realm, clientID, clientSecret)
	graphQLClient := graphql.NewClient(serviceClient, assetManagerURL)

	assetRepository := graphql.NewAssetRepository(graphQLClient, circuitBreaker)
	bucketRepository := graphql.NewBucketRepository(graphQLClient, circuitBreaker)

	assetService := assetapp.NewApplicationService(assetRepository, cacheService, circuitBreaker)
	bucketService := bucketapp.NewApplicationService(bucketRepository, cacheService, circuitBreaker)

	handler := httphandler.NewHandler(assetService, bucketService, cfg)
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
