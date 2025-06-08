package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	// Parse request body
	var req UploadRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil || strings.TrimSpace(req.FileName) == "" {
		log.Printf("Invalid request body: %v | raw body: %s", err, event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "fileName is required in the request body"})
	}

	// Validate environment variables
	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		log.Println("Missing BUCKET_NAME env variable")
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server configuration error: missing bucket name"})
	}

	region := os.Getenv("BUCKET_REGION")
	if region == "" {
		region = "eu-north-1"
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to initialize AWS session"})
	}

	svc := s3.New(sess)

	// Create pre-signed URL for PUT
	reqObj, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(req.FileName),
	})
	url, err := reqObj.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Failed to generate pre-signed URL: %v", err)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate pre-signed URL"})
	}

	return respondJSON(http.StatusOK, UploadResponse{URL: url})
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
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*", //TODO: make this a variable later
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}