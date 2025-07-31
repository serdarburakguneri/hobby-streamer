package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type DeleteRequest struct {
	AssetID string `json:"assetId"`
	Folder  string `json:"folder"`
}

type DeleteResponse struct {
	Message string   `json:"message"`
	Folder  string   `json:"folder"`
	Deleted int      `json:"deleted"`
	Errors  []string `json:"errors,omitempty"`
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

	log := logger.WithService("delete-files")
	if trackingID != "" {
		log = log.WithTrackingID(trackingID)
	}
	log = log.WithContext(ctx)

	if event.Body == "" {
		log.Error("Empty request body")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Request body is required",
			Type:    "validation",
		})
	}

	var req DeleteRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		log.WithError(err).Error("Invalid request body", "raw_body", event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request body format",
			Type:    "validation",
		})
	}

	if req.AssetID == "" {
		log.Error("Missing assetId in request")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "assetId is required",
			Type:    "validation",
		})
	}

	if req.Folder == "" {
		log.Error("Missing folder in request", "asset_id", req.AssetID)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "folder is required",
			Type:    "validation",
		})
	}

	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsConfig := &aws.Config{Region: aws.String(region)}
	if endpoint != "" {
		awsConfig.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.WithError(err).Error("Failed to create AWS session", "asset_id", req.AssetID)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to initialize AWS session",
			Type:    "internal",
		})
	}

	svc := s3.New(sess)

	parts := strings.SplitN(req.Folder, "/", 2)
	if len(parts) != 2 {
		log.Error("Invalid folder format", "folder", req.Folder, "asset_id", req.AssetID)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid folder format. Expected: bucket/prefix",
			Type:    "validation",
		})
	}

	bucket := parts[0]
	prefix := parts[1] + "/"

	log.Info("Starting folder deletion", "asset_id", req.AssetID, "bucket", bucket, "prefix", prefix)

	var allObjects []*s3.ObjectIdentifier
	var errors []string

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	err = svc.ListObjectsV2PagesWithContext(ctx, listInput, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			allObjects = append(allObjects, &s3.ObjectIdentifier{
				Key: obj.Key,
			})
		}
		return !lastPage
	})

	if err != nil {
		log.WithError(err).Error("Failed to list objects in folder", "bucket", bucket, "prefix", prefix, "asset_id", req.AssetID)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to list objects in folder",
			Type:    "internal",
		})
	}

	if len(allObjects) == 0 {
		log.Info("No objects found in folder", "asset_id", req.AssetID, "bucket", bucket, "prefix", prefix)
		return respondJSON(http.StatusOK, DeleteResponse{
			Message: "No objects found in folder",
			Folder:  req.Folder,
			Deleted: 0,
		})
	}

	batchSize := 1000
	totalDeleted := 0

	for i := 0; i < len(allObjects); i += batchSize {
		end := i + batchSize
		if end > len(allObjects) {
			end = len(allObjects)
		}

		batch := allObjects[i:end]

		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &s3.Delete{
				Objects: batch,
				Quiet:   aws.Bool(false),
			},
		}

		result, err := svc.DeleteObjectsWithContext(ctx, input)
		if err != nil {
			log.WithError(err).Error("Failed to delete batch of objects", "bucket", bucket, "batch_size", len(batch), "asset_id", req.AssetID)
			errors = append(errors, fmt.Sprintf("Failed to delete batch: %v", err))
			continue
		}

		totalDeleted += len(result.Deleted)

		for _, err := range result.Errors {
			errorMsg := fmt.Sprintf("Failed to delete %s: %s", *err.Key, *err.Message)
			errors = append(errors, errorMsg)
			log.Error("Delete error", "bucket", bucket, "key", *err.Key, "error", *err.Message, "asset_id", req.AssetID)
		}
	}

	log.Info("Folder deletion completed", "asset_id", req.AssetID, "deleted_count", totalDeleted, "error_count", len(errors))

	response := DeleteResponse{
		Message: fmt.Sprintf("Deleted %d objects from folder %s", totalDeleted, req.Folder),
		Folder:  req.Folder,
		Deleted: totalDeleted,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	apiResponse, err := respondJSON(http.StatusOK, response)
	if err != nil {
		return apiResponse, err
	}

	if trackingID != "" {
		apiResponse.Headers["X-Tracking-ID"] = trackingID
	}

	return apiResponse, nil
}

func respondJSON(status int, payload interface{}) (events.APIGatewayProxyResponse, error) {
	log := logger.WithService("delete-files")

	body, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Failed to marshal JSON response")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	logger.Init(slog.LevelInfo, "json")
	lambda.Start(handler)
}
