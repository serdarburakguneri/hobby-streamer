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
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

func main() {
	logLevel := getLogLevel()
	logFormat := getEnv("LOG_FORMAT", "text")
	logger.Init(logLevel, logFormat)
	log := logger.WithService("asset-manager")
	log.Info("Starting asset-manager service")

	neo4jURI := getEnv("NEO4J_URI", "bolt://localhost:7687")
	neo4jUsername := getEnv("NEO4J_USERNAME", "neo4j")
	neo4jPassword := getEnv("NEO4J_PASSWORD", "password")
	port := getEnv("PORT", "8080")

	driver, err := neo4j.NewDriver(neo4jURI, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))
	if err != nil {
		log.WithError(err).Error("Failed to create Neo4j driver")
		os.Exit(1)
	}
	defer driver.Close()

	if err := driver.VerifyConnectivity(); err != nil {
		log.WithError(err).Error("Failed to connect to Neo4j")
		os.Exit(1)
	}

	assetRepo := asset.NewRepository(driver)
	bucketRepo := bucket.NewRepository(driver)

	sqsQueueURL := getEnv("TRANSCODER_QUEUE_URL", "")
	var assetService *asset.Service
	if sqsQueueURL != "" {
		sqsProducer, err := sqs.NewProducer(context.Background(), sqsQueueURL)
		if err != nil {
			log.WithError(err).Error("Failed to create SQS producer")
			os.Exit(1)
		}
		assetService = asset.NewServiceWithSQS(assetRepo, sqsProducer)
		log.Info("Asset service initialized with SQS producer", "queue_url", sqsQueueURL)
	} else {
		assetService = asset.NewService(assetRepo)
		log.Info("Asset service initialized without SQS producer")
	}

	bucketService := bucket.NewService(bucketRepo)

	analyzeQueueURL := getEnv("ANALYZE_QUEUE_URL", "")
	var consumerRegistry *sqs.ConsumerRegistry
	if analyzeQueueURL != "" {
		consumerRegistry = sqs.NewConsumerRegistry()

		messageRouter := consumer.NewMessageRouter(assetService)
		consumerRegistry.Register(analyzeQueueURL, messageRouter.HandleMessage)

		go func() {
			if err := consumerRegistry.Start(context.Background()); err != nil {
				log.WithError(err).Error("Failed to start consumer registry")
			}
		}()
		log.Info("Message router initialized", "queue_url", analyzeQueueURL, "supported_messages", []string{"analyze-completed", "transcode-hls-completed", "transcode-dash-completed"})
	} else {
		log.Info("Message router not initialized - no queue URL provided")
	}

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

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			keycloakURL := getEnv("KEYCLOAK_URL", "http://localhost:8080")
			realm := getEnv("KEYCLOAK_REALM", "hobby-realm")
			clientID := getEnv("KEYCLOAK_CLIENT_ID", "asset-manager")

			log := logger.WithService("auth-middleware")
			log.Debug("Validating token", "keycloakURL", keycloakURL, "realm", realm, "clientID", clientID)

			validator := auth.NewKeycloakValidator(keycloakURL, realm, clientID)
			user, err := validator.ValidateToken(r.Context(), token)
			if err != nil {
				log.WithError(err).Debug("Regular token validation failed, trying service token")

				serviceValidator := auth.NewServiceTokenValidator(keycloakURL, realm, clientID)
				serviceUser, serviceErr := serviceValidator.ValidateServiceToken(r.Context(), token)
				if serviceErr != nil {
					log.WithError(serviceErr).Error("Service token validation also failed")
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}

				if !serviceValidator.IsServiceToken(serviceUser) {
					log.Error("Service token is not from streaming-api service")
					http.Error(w, "Invalid service token", http.StatusUnauthorized)
					return
				}

				log.Debug("Service token validated successfully", "service", serviceUser.ClientID, "roles", serviceUser.Roles)
				ctx := context.WithValue(r.Context(), "service_user", serviceUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Debug("Token validated successfully", "user", user.Username, "roles", user.Roles)
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	router.Use(corsMiddleware)
	router.Use(loggerMiddleware(log))
	router.Use(authMiddleware)

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
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Info("Starting server", "port", port)
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getLogLevel() slog.Level {
	level := getEnv("LOG_LEVEL", "info")
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
