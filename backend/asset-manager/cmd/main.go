package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	neo4jdriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	appasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	appbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket"
	domainasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	neo4jrepo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/infrastructure/neo4j"
	sqsevents "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/infrastructure/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/interfaces/graphql"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
	pkgsqs "github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

const defaultPort = "8080"

func setupLogger(cfg *config.BaseConfig) {
	if cfg.Log.Async.Enabled {
		logger.InitAsync(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format, cfg.Log.Async.BufferSize)
	} else {
		logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	}
}

func setupNeo4j(dynamicCfg *config.DynamicConfig, secretsManager *config.SecretsManager) neo4jdriver.Driver {
	neo4jURI := dynamicCfg.GetStringFromComponent("neo4j", "uri")
	neo4jUsername := dynamicCfg.GetStringFromComponent("neo4j", "username")
	neo4jPassword := secretsManager.Get("neo4j_password")
	if neo4jPassword == "" {
		logger.Get().Error("Failed to create Neo4j driver: Neo4j password is empty")
		os.Exit(1)
	}

	neo4jDriver, err := neo4jdriver.NewDriver(neo4jURI, neo4jdriver.BasicAuth(neo4jUsername, neo4jPassword, ""))
	if err != nil {
		logger.Get().Error("Failed to create Neo4j driver", "error", err)
		os.Exit(1)
	}

	if err := neo4jDriver.VerifyConnectivity(); err != nil {
		logger.Get().Error("Failed to connect to Neo4j", "error", err)
		os.Exit(1)
	}

	return neo4jDriver
}

func setupSQS(ctx context.Context, dynamicCfg *config.DynamicConfig) (*pkgsqs.Producer, *pkgsqs.Producer) {
	jobQueueURL := dynamicCfg.GetStringFromComponent("sqs", "job_queue_url")
	domainEventQueueURL := "http://localhost:4566/000000000000/asset-events"

	domainEventProducer, err := pkgsqs.NewProducer(ctx, domainEventQueueURL)
	if err != nil {
		logger.Get().Error("Failed to create domain event producer", "error", err)
		os.Exit(1)
	}

	jobProducer, err := pkgsqs.NewProducer(ctx, jobQueueURL)
	if err != nil {
		logger.Get().Error("Failed to create job producer", "error", err)
		os.Exit(1)
	}

	return domainEventProducer, jobProducer
}

func setupServices(neo4jDriver neo4jdriver.Driver, domainEventProducer, jobProducer *pkgsqs.Producer) (*appasset.ApplicationService, *appbucket.ApplicationService) {
	assetRepo := neo4jrepo.NewAssetRepository(neo4jDriver)
	bucketRepo := neo4jrepo.NewBucketRepository(neo4jDriver)

	eventPublisher := sqsevents.NewEventPublisherWithJobProducer(domainEventProducer, jobProducer)

	assetDomainService := domainasset.NewDomainService(assetRepo)
	assetPublishingService := domainasset.NewPublishingService(assetDomainService)
	bucketDomainService := domainbucket.NewDomainService(bucketRepo)

	assetAppService := appasset.NewApplicationService(assetRepo, assetDomainService, assetPublishingService, eventPublisher)
	bucketAppService := appbucket.NewApplicationService(bucketRepo, bucketDomainService, eventPublisher)

	return assetAppService, bucketAppService
}

func setupGraphQL(assetAppService *appasset.ApplicationService, bucketAppService *appbucket.ApplicationService, cfg *config.BaseConfig) *handler.Server {
	resolver := graphql.NewResolver(assetAppService, bucketAppService)
	schema := graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver})
	gqlHandler := handler.New(schema)

	gqlHandler.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				for _, allowed := range cfg.Security.CORS.AllowedOrigins {
					if origin == allowed {
						return true
					}
				}
				return false
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	gqlHandler.AddTransport(&transport.Options{})
	gqlHandler.AddTransport(&transport.GET{})
	gqlHandler.AddTransport(&transport.POST{})
	gqlHandler.AddTransport(&transport.MultipartForm{})
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.FixedComplexityLimit(1000))

	return gqlHandler
}

func setupAuth(dynamicCfg *config.DynamicConfig) func(http.HandlerFunc) http.HandlerFunc {
	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	keycloakRealm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	keycloakClientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")

	authValidator := auth.NewKeycloakValidator(keycloakURL, keycloakRealm, keycloakClientID)
	authMiddleware := auth.NewAuthMiddleware(authValidator)
	return authMiddleware.RequireUserAuth().RequireServiceAuth().Build()
}

func setupRouter(gqlHandler *handler.Server, authHandlerFunc func(http.HandlerFunc) http.HandlerFunc) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/graphql", authHandlerFunc(gqlHandler.ServeHTTP)).Methods("GET", "POST", "OPTIONS")
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return router
}

func setupMiddleware(router *mux.Router, cfg *config.BaseConfig) http.Handler {
	corsMiddleware := security.CORSMiddleware(
		cfg.Security.CORS.AllowedOrigins,
		cfg.Security.CORS.AllowedMethods,
		cfg.Security.CORS.AllowedHeaders,
	)

	handler := corsMiddleware(router)
	handler = security.RateLimitMiddleware(cfg.Security.RateLimit.Requests, cfg.Security.RateLimit.Window)(handler)
	handler = security.SecurityHeadersMiddleware()(handler)

	return handler
}

func setupServer(handler http.Handler, cfg *config.BaseConfig) *http.Server {
	port := cfg.Server.Port
	if port == "" {
		port = defaultPort
	}

	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func main() {
	ctx := context.Background()

	configManager, err := config.NewManager("asset-manager")
	if err != nil {
		logger.Get().Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}
	defer configManager.Close()

	secretsManager := config.NewSecretsManager()
	secretsManager.LoadFromEnvironment()

	cfg := configManager.GetConfig()
	dynamicCfg := configManager.GetDynamicConfig()

	setupLogger(cfg)
	slog := logger.WithService(cfg.Service)
	slog.Info("Starting asset-manager service", "environment", cfg.Environment)

	neo4jDriver := setupNeo4j(dynamicCfg, secretsManager)
	defer neo4jDriver.Close()

	domainEventProducer, jobProducer := setupSQS(ctx, dynamicCfg)
	assetAppService, bucketAppService := setupServices(neo4jDriver, domainEventProducer, jobProducer)

	eventConsumer := sqsevents.NewEventConsumer(assetAppService)
	completionQueueURL := dynamicCfg.GetStringFromComponent("sqs", "completion_queue_url")
	eventConsumer.RegisterCompletionQueue(completionQueueURL)

	go func() {
		if err := eventConsumer.Start(ctx); err != nil {
			slog.WithError(err).Error("Failed to start event consumer")
		}
	}()

	gqlHandler := setupGraphQL(assetAppService, bucketAppService, cfg)
	authHandlerFunc := setupAuth(dynamicCfg)
	router := setupRouter(gqlHandler, authHandlerFunc)
	handler := setupMiddleware(router, cfg)
	server := setupServer(handler, cfg)

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
