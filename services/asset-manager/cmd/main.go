package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

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

	r := router.NewRouter(assetService, bucketService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting asset-manager on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
