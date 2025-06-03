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
	limit := 10
	if val, ok := req.QueryStringParameters["limit"]; ok {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}

	var lastKey map[string]map[string]string
	if token, ok := req.QueryStringParameters["nextKey"]; ok && token != "" {
		parsed, err := asset.DecodeLastEvaluatedKey(token)
		if err != nil {
			log.Printf("Invalid nextKey: %v", err)
			return response(http.StatusBadRequest, "Invalid nextKey")
		}
		lastKey = parsed
	}

	scanKey, err := asset.ToDynamoKey(lastKey)
	if err != nil {
		log.Printf("Failed to convert scan key: %v", err)
		return response(http.StatusBadRequest, "Invalid scan key")
	}

	page, err := repo.ListAssets(ctx, limit, scanKey)
	if err != nil {
		log.Printf("ListAssets failed: %v", err)
		return response(http.StatusInternalServerError, "Failed to list assets")
	}

	resp := asset.BuildPaginatedResponse(page)
	return jsonResponse(http.StatusOK, resp)
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