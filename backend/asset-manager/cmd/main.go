package main

import (
	"context"
	"log/slog"
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
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/graph"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/consumer"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

func main() {
	configManager, err := config.NewManager("asset-manager")
	if err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}
	defer configManager.Close()

	secretsManager := config.NewSecretsManager()
	secretsManager.LoadFromEnvironment()

	cfg := configManager.GetConfig()
	dynamicCfg := configManager.GetDynamicConfig()

	logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	log := logger.WithService(cfg.Service)
	log.Info("Starting asset-manager service", "environment", cfg.Environment)

	uri := dynamicCfg.GetStringFromComponent("neo4j", "uri")
	username := dynamicCfg.GetStringFromComponent("neo4j", "username")
	password := secretsManager.Get("neo4j_password")

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.WithError(err).Error("Failed to create Neo4j driver")
		os.Exit(1)
	}
	defer driver.Close()

	if err := driver.VerifyConnectivity(); err != nil {
		log.WithError(err).Error("Failed to connect to Neo4j")
		os.Exit(1)
	}
	log.Info("Neo4j connection established", "uri", uri)

	assetRepo := asset.NewRepository(driver)
	bucketRepo := bucket.NewRepository(driver)

	transcoderQueueURL := dynamicCfg.GetStringFromComponent("sqs", "transcoder_queue_url")
	sqsProducer, err := sqs.NewProducer(context.Background(), transcoderQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create SQS producer")
		os.Exit(1)
	}
	assetService := asset.NewServiceWithSQS(assetRepo, sqsProducer, dynamicCfg)
	log.Info("Asset service initialized with SQS producer", "queue_url", transcoderQueueURL)

	bucketService := bucket.NewService(bucketRepo)

	analyzeQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_queue_url")
	consumerRegistry := sqs.NewConsumerRegistry()

	messageRouter := consumer.NewMessageRouter(assetService)
	consumerRegistry.Register(analyzeQueueURL, messageRouter.HandleMessage)

	go func() {
		if err := consumerRegistry.Start(context.Background()); err != nil {
			log.WithError(err).Error("Failed to start consumer registry")
		}
	}()
	log.Info("Message router initialized", "queue_url", analyzeQueueURL, "supported_messages", []string{"analyze-completed", "transcode-hls-completed", "transcode-dash-completed"})

	router := mux.NewRouter()

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	loggerMiddleware := func(log *logger.Logger) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
				next.ServeHTTP(w, r)
			})
		}
	}

	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	keycloakRealm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	keycloakClientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")

	keycloakValidator := auth.NewKeycloakValidator(keycloakURL, keycloakRealm, keycloakClientID)
	authMiddleware := auth.NewAuthMiddleware(keycloakValidator)

	router.Use(corsMiddleware)
	router.Use(loggerMiddleware(log))
	router.Use(func(next http.Handler) http.Handler {
		return authMiddleware.RequireUserAuth().RequireServiceAuth().Build()(next.ServeHTTP)
	})

	resolver := graph.NewResolver(assetService, bucketService)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})
	gqlHandler := handler.New(schema)

	gqlHandler.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	gqlHandler.AddTransport(&transport.Options{})
	gqlHandler.AddTransport(&transport.GET{})
	gqlHandler.AddTransport(&transport.POST{})
	gqlHandler.AddTransport(&transport.MultipartForm{})
	gqlHandler.Use(extension.Introspection{})

	router.Handle("/graphql", gqlHandler).Methods("GET", "POST", "OPTIONS")
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("Starting server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
