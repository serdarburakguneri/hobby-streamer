package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/serdarburakguneri/hobby-streamer/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/bucket"
	router "github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/http"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	assetRepo := asset.NewRepository("asset", dynamoClient)
	assetService := asset.NewService(assetRepo)

	bucketRepo := bucket.NewRepository("bucket", dynamoClient)
	bucketService := bucket.NewService(bucketRepo)

	keycloakURL := getEnv("KEYCLOAK_URL", "http://localhost:8080")
	realm := getEnv("KEYCLOAK_REALM", "hobby")
	clientID := getEnv("KEYCLOAK_CLIENT_ID", "asset-manager")

	validator := auth.NewKeycloakValidator(keycloakURL, realm, clientID)
	authMiddleware := auth.NewAuthMiddleware(validator)

	r := router.NewRouter(assetService, bucketService, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting asset-manager on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
