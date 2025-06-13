package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/http"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/asset"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo := asset.NewRepository("asset", dynamoClient)

	router := http.NewRouter(repo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting asset-manager on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}