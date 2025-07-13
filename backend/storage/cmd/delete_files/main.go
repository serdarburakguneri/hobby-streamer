package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request body
	var req DeleteRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		log.Printf("Invalid request body: %v | raw body: %s", err, event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
	}

	if req.AssetID == "" {
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "assetId is required"})
	}

	if len(req.Files) == 0 {
		return respondJSON(http.StatusOK, DeleteResponse{
			Message: "No files to delete",
			Deleted: []S3File{},
		})
	}

	// Support custom endpoint for LocalStack
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsConfig := &aws.Config{Region: aws.String(region)}
	if endpoint != "" {
		awsConfig.Endpoint = aws.String(endpoint)
	}

	// Create AWS session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to initialize AWS session"})
	}

	svc := s3.New(sess)

	// Group files by bucket for efficient deletion
	bucketFiles := make(map[string][]string)
	for _, file := range req.Files {
		bucketFiles[file.Bucket] = append(bucketFiles[file.Bucket], file.Key)
	}

	var deleted []S3File
	var errors []string

	// Delete files from each bucket
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
			log.Printf("Failed to delete objects from bucket %s: %v", bucket, err)
			errors = append(errors, fmt.Sprintf("Failed to delete from bucket %s: %v", bucket, err))
			continue
		}

		// Add successfully deleted files
		for _, deletedObj := range result.Deleted {
			deleted = append(deleted, S3File{
				Bucket: bucket,
				Key:    *deletedObj.Key,
			})
		}

		// Log any errors
		for _, err := range result.Errors {
			errorMsg := fmt.Sprintf("Failed to delete %s from %s: %s", *err.Key, bucket, *err.Message)
			errors = append(errors, errorMsg)
			log.Printf("Delete error: %s", errorMsg)
		}
	}

	response := DeleteResponse{
		Message: fmt.Sprintf("Deleted %d files for asset %s", len(deleted), req.AssetID),
		Deleted: deleted,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	return respondJSON(http.StatusOK, response)
}

func respondJSON(status int, payload interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", err)
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
	lambda.Start(handler)
}
