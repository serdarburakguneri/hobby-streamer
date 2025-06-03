package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"asset-manager/internal/asset"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var repo *asset.Repository

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo = asset.NewRepository("asset", dynamoClient)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var a asset.Asset

	if err := json.Unmarshal([]byte(req.Body), &a); err != nil {
		log.Printf("Invalid input: %v", err)
		return response(http.StatusBadRequest, "Invalid request body")
	}

	if err := repo.SaveAsset(ctx, &a); err != nil {
		log.Printf("Failed to save asset: %v", err)
		return response(http.StatusInternalServerError, "Could not save asset")
	}

	respBody, _ := json.Marshal(a)
	return response(http.StatusOK, string(respBody))
}

func response(status int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}