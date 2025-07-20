package main

import (
	"context"
	"encoding/json"
	"fmt"
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

type ImageUploadRequest struct {
	FileName  string `json:"fileName"`
	AssetID   string `json:"assetId"`
	ImageType string `json:"imageType"`
}

type UploadResponse struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	trackingID := ""
	if trackingHeader, exists := event.Headers["X-Tracking-ID"]; exists {
		trackingID = trackingHeader
	} else if trackingHeader, exists := event.Headers["x-tracking-id"]; exists {
		trackingID = trackingHeader
	}

	log := logger.WithService("generate-image-upload-url")
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

	if event.Body == "" {
		log.Error("Empty request body")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Request body is required",
			Type:    "validation",
		})
	}

	var req ImageUploadRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		log.WithError(err).Error("Invalid request body", "raw_body", event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request body format",
			Type:    "validation",
		})
	}

	if strings.TrimSpace(req.FileName) == "" {
		log.Error("Missing or empty fileName in request")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "fileName is required and cannot be empty",
			Type:    "validation",
		})
	}

	if strings.TrimSpace(req.AssetID) == "" {
		log.Error("Missing assetId in request")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "assetId is required",
			Type:    "validation",
		})
	}

	if strings.TrimSpace(req.ImageType) == "" {
		log.Error("Missing imageType in request")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "imageType is required",
			Type:    "validation",
		})
	}

	if strings.Contains(req.FileName, "..") {
		log.Error("Invalid fileName - path traversal attempt", "file_name", req.FileName)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "fileName contains invalid characters",
			Type:    "validation",
		})
	}

	if strings.Contains(req.FileName, "/") {
		log.Error("Invalid fileName - contains path separators", "file_name", req.FileName)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "fileName cannot contain path separators",
			Type:    "validation",
		})
	}

	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		bucket = "content-east"
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
		log.WithError(err).Error("Failed to create AWS session", "file_name", req.FileName, "asset_id", req.AssetID)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to initialize AWS session",
			Type:    "internal",
		})
	}

	svc := s3.New(sess)

	s3Key := fmt.Sprintf("%s/images/%s/%s", req.AssetID, strings.ToLower(req.ImageType), req.FileName)

	reqObj, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3Key),
	})
	url, err := reqObj.Presign(15 * time.Minute)
	if err != nil {
		log.WithError(err).Error("Failed to generate pre-signed URL", "file_name", req.FileName, "asset_id", req.AssetID, "s3_key", s3Key, "bucket", bucket)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to generate pre-signed URL",
			Type:    "external",
		})
	}

	log.Info("Generated presigned URL for image", "bucket", bucket, "key", s3Key, "file_name", req.FileName, "asset_id", req.AssetID, "image_type", req.ImageType, "expires_in_minutes", 15)

	response, err := respondJSON(http.StatusOK, UploadResponse{URL: url})
	if err != nil {
		return response, err
	}

	if trackingID != "" {
		response.Headers["X-Tracking-ID"] = trackingID
	}

	return response, nil
}

func respondJSON(status int, payload interface{}) (events.APIGatewayProxyResponse, error) {
	log := logger.WithService("generate-image-upload-url")

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
			"Access-Control-Allow-Origin":      "http://localhost:8081",
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
