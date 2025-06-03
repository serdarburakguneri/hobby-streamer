package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

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
	idStr, ok := req.PathParameters["id"]
	if !ok || idStr == "" {
		return response(http.StatusBadRequest, "Missing asset ID in path")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid asset ID: %v", err)
		return response(http.StatusBadRequest, "Asset ID must be a number")
	}

	a, err := repo.GetAssetByID(ctx, id)
	if err != nil {
		log.Printf("GetAssetByID failed: %v", err)
		return response(http.StatusNotFound, "Asset not found")
	}

	return jsonResponse(http.StatusOK, a)
}

func response(status int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       `"` + message + `"`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func jsonResponse(status int, data interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := asset.SafeJSON(data)
	if err != nil {
		log.Printf("Failed to serialize response: %v", err)
		return response(http.StatusInternalServerError, "Internal error")
	}

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