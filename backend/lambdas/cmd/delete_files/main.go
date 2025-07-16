package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type DeleteRequest struct {
	AssetID string   `json:"assetId"`
	Files   []S3File `json:"files"`
}

type S3File struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type DeleteResponse struct {
	Message string   `json:"message"`
	Deleted []S3File `json:"deleted"`
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

	if req.Files == nil {
		log.Error("Files array is nil", "asset_id", req.AssetID)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "files array is required",
			Type:    "validation",
		})
	}

	if len(req.Files) == 0 {
		log.Info("No files to delete", "asset_id", req.AssetID)
		return respondJSON(http.StatusOK, DeleteResponse{
			Message: "No files to delete",
			Deleted: []S3File{},
		})
	}

	for i, file := range req.Files {
		if file.Bucket == "" {
			log.Error("Missing bucket in file", "asset_id", req.AssetID, "file_index", i)
			return respondJSON(http.StatusBadRequest, ErrorResponse{
				Message: fmt.Sprintf("bucket is required for file at index %d", i),
				Type:    "validation",
			})
		}
		if file.Key == "" {
			log.Error("Missing key in file", "asset_id", req.AssetID, "file_index", i, "bucket", file.Bucket)
			return respondJSON(http.StatusBadRequest, ErrorResponse{
				Message: fmt.Sprintf("key is required for file at index %d", i),
				Type:    "validation",
			})
		}
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

	bucketFiles := make(map[string][]string)
	for _, file := range req.Files {
		bucketFiles[file.Bucket] = append(bucketFiles[file.Bucket], file.Key)
	}

	var deleted []S3File
	var errors []string

	for bucket, keys := range bucketFiles {
		if len(keys) == 0 {
			continue
		}

		var objects []*s3.ObjectIdentifier
		for _, key := range keys {
			objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(key)})
		}

		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &s3.Delete{
				Objects: objects,
				Quiet:   aws.Bool(false),
			},
		}

		result, err := svc.DeleteObjectsWithContext(ctx, input)
		if err != nil {
			log.WithError(err).Error("Failed to delete objects from bucket", "bucket", bucket, "file_count", len(keys), "asset_id", req.AssetID)
			errors = append(errors, fmt.Sprintf("Failed to delete from bucket %s: %v", bucket, err))
			continue
		}

		for _, deletedObj := range result.Deleted {
			deleted = append(deleted, S3File{
				Bucket: bucket,
				Key:    *deletedObj.Key,
			})
		}

		for _, err := range result.Errors {
			errorMsg := fmt.Sprintf("Failed to delete %s from %s: %s", *err.Key, bucket, *err.Message)
			errors = append(errors, errorMsg)
			log.Error("Delete error", "bucket", bucket, "key", *err.Key, "error", *err.Message, "asset_id", req.AssetID)
		}
	}

	log.Info("File deletion completed", "asset_id", req.AssetID, "deleted_count", len(deleted), "error_count", len(errors))

	response := DeleteResponse{
		Message: fmt.Sprintf("Deleted %d files for asset %s", len(deleted), req.AssetID),
		Deleted: deleted,
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
