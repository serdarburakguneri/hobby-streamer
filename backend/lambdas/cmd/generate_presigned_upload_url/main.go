package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type UploadRequest struct {
	FileName string `json:"fileName"`
}

type UploadResponse struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract tracking ID from headers
	trackingID := ""
	if trackingHeader, exists := event.Headers["X-Tracking-ID"]; exists {
		trackingID = trackingHeader
	} else if trackingHeader, exists := event.Headers["x-tracking-id"]; exists {
		trackingID = trackingHeader
	}
	
	// Create logger with tracking ID
	log := logger.WithService("generate-presigned-url")
	if trackingID != "" {
		log = log.WithTrackingID(trackingID)
	}
	log = log.WithContext(ctx)

	if event.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With",
				"Access-Control-Allow-Credentials": "true",
			},
		}, nil
	}

	var req UploadRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil || strings.TrimSpace(req.FileName) == "" {
		log.WithError(err).Error("Invalid request body", "raw_body", event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "fileName is required in the request body"})
	}

	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		log.Error("Missing BUCKET_NAME env variable")
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server configuration error: missing bucket name"})
	}

	region := os.Getenv("BUCKET_REGION")
	if region == "" {
		region = "us-east-1"
	}

	endpoint := os.Getenv("AWS_ENDPOINT")
	awsConfig := &aws.Config{Region: aws.String(region)}
	if endpoint != "" {
		awsConfig.Endpoint = aws.String(endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.WithError(err).Error("Failed to create AWS session")
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to initialize AWS session"})
	}

	svc := s3.New(sess)

	reqObj, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(req.FileName),
	})
	url, err := reqObj.Presign(15 * time.Minute)
	if err != nil {
		log.WithError(err).Error("Failed to generate pre-signed URL")
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate pre-signed URL"})
	}

	log.Info("Generated presigned URL", "bucket", bucket, "key", req.FileName, "expires_in_minutes", 15)

	response, err := respondJSON(http.StatusOK, UploadResponse{URL: url})
	if err != nil {
		return response, err
	}
	
	// Add tracking ID to response headers
	if trackingID != "" {
		response.Headers["X-Tracking-ID"] = trackingID
	}
	
	return response, nil
}

func respondJSON(status int, payload interface{}) (events.APIGatewayProxyResponse, error) {
	log := logger.WithService("generate-presigned-url")

	body, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Failed to marshal JSON response")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func main() {
	logger.Init(slog.LevelInfo, "json")
	lambda.Start(handler)
}
