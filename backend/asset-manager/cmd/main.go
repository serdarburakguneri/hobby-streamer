package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/graph"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

func main() {
	log := logger.Get().WithService("asset-manager")
	log.Info("Starting asset-manager service")

	neo4jURI := getEnv("NEO4J_URI", "bolt://localhost:7687")
	neo4jUsername := getEnv("NEO4J_USERNAME", "neo4j")
	neo4jPassword := getEnv("NEO4J_PASSWORD", "password")
	port := getEnv("PORT", "8080")
	sqsQueueURL := getEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/transcoder-jobs")
	statusQueueURL := getEnv("STATUS_QUEUE_URL", "http://localhost:4566/000000000000/status-updates")

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

	var assetService asset.AssetService
	if sqsQueueURL != "" {
		sqsProducer, err := sqs.NewProducer(context.Background(), sqsQueueURL)
		if err != nil {
			log.WithError(err).Error("Failed to create SQS producer, falling back to service without SQS")
			assetService = asset.NewService(assetRepo)
		} else {
			log.Info("SQS producer initialized successfully", "queue_url", sqsQueueURL)
			assetService = asset.NewServiceWithSQS(assetRepo, sqsProducer)
		}
	} else {
		assetService = asset.NewService(assetRepo)
	}

	if statusQueueURL != "" {
		statusConsumer, err := sqs.NewConsumer(context.Background(), statusQueueURL)
		if err != nil {
			log.WithError(err).Error("Failed to create status queue consumer, continuing without status updates")
		} else {
			log.Info("Status queue consumer initialized successfully", "status_queue_url", statusQueueURL)
			go func() {
				statusConsumer.Start(context.Background(), func(msg sqs.Message) error {
					var payload map[string]interface{}
					if err := json.Unmarshal(msg.Payload, &payload); err != nil {
						log.WithError(err).Error("Failed to unmarshal status message payload")
						return err
					}
					return assetService.HandleStatusUpdateMessage(context.Background(), msg.Type, payload)
				})
			}()
		}
	}

	bucketService := bucket.NewService(bucketRepo)

	router := mux.NewRouter()

	router.Use(corsMiddleware)
	router.Use(loggerMiddleware(log))
	router.Use(authMiddleware)

	resolver := graph.NewResolver(assetService, bucketService)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})

	srv := handler.NewDefaultServer(schema)

	// Use simple http.Handle for GraphQL routes
	http.Handle("/graphql", srv)
	http.Handle("/graphql/", srv)

	// Mount the http.Handle routes on the main router
	router.PathPrefix("/graphql").Handler(http.DefaultServeMux)

	if getEnv("ENV", "development") == "development" {
		playgroundHandler := playground.Handler("GraphQL", "/graphql")
		router.Handle("/playground", playgroundHandler)
		log.Info("GraphQL playground available at /playground")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
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

		log := logger.Get().WithService("auth-middleware")
		log.Debug("Validating token", "keycloakURL", keycloakURL, "realm", realm, "clientID", clientID)

		validator := auth.NewKeycloakValidator(keycloakURL, realm, clientID)
		user, err := validator.ValidateToken(r.Context(), token)
		if err != nil {
			log.WithError(err).Error("Token validation failed")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Debug("Token validated successfully", "user", user.Username, "roles", user.Roles)
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggerMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			log.Info("HTTP Request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
				"user_agent", r.UserAgent(),
			)
		})
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
