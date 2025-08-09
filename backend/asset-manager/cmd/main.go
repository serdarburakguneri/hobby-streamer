package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bootstrap"
	outbox "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/infrastructure/outbox"
	bootstrap_events "github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfgManager, secretsManager, cfg, dynamicCfg, err := bootstrap.InitConfig("asset-manager")
	if err != nil {
		logger.Get().Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}
	defer cfgManager.Close()

	bootstrap.InitLogger(cfg)
	slog := logger.WithService(cfg.Service)
	slog.Info("Starting asset-manager service", "environment", cfg.Environment)

	neo4jDriver := bootstrap.InitNeo4j(dynamicCfg, secretsManager)
	defer neo4jDriver.Close()

	cdnService := bootstrap.InitCDNService(dynamicCfg)

	domainProducer, jobProducer := bootstrap.InitKafkaProducers(ctx, dynamicCfg)

	assetCmdService, assetQryService, bucketCmdService, bucketQryService, pipelineService, _ := bootstrap.InitServices(neo4jDriver)

	assetEventConsumer := bootstrap.InitKafkaConsumer(ctx, assetCmdService, assetQryService, domainProducer, cdnService, pipelineService, dynamicCfg, neo4jDriver)
	defer assetEventConsumer.Stop()

	var gqlPublisher interface {
		Publish(ctx context.Context, topic string, ev *bootstrap_events.Event) error
	}
	if dynamicCfg.GetBaseConfig().Features.EnableOutbox {
		store := outbox.NewNeo4jStore(neo4jDriver)
		gqlPublisher = outbox.NewPublisher(store, jobProducer)
		dispatcher := outbox.NewDispatcher(store, jobProducer)
		dispatcher.Start(ctx)
	} else {
		gqlPublisher = jobProducer
	}
	gqlHandler := bootstrap.InitGraphQL(assetCmdService, assetQryService, bucketCmdService, bucketQryService, cdnService, pipelineService, gqlPublisher, cfg)
	authHandlerFunc := bootstrap.InitAuth(dynamicCfg)
	router := bootstrap.InitRouter(gqlHandler, authHandlerFunc)
	handler := bootstrap.InitMiddleware(router, cfg)
	server := bootstrap.InitServer(handler, cfg)

	go func() {
		slog.Info("Starting HTTP server", "port", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.WithError(err).Error("Failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.WithError(err).Error("Server forced to shutdown")
		os.Exit(1)
	}

	slog.Info("Server exited")
}
