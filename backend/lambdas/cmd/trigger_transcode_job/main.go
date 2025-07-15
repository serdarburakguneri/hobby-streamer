package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type TranscodeRequest struct {
	AssetID        string `json:"assetId"`
	VideoType      string `json:"videoType"`
	Format         string `json:"format"`
	Input          string `json:"input,omitempty"`
	SourceFileName string `json:"sourceFileName,omitempty"`
}

type TranscodeResponse struct {
	Message string `json:"message"`
	JobType string `json:"jobType"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	if event.HTTPMethod != "POST" {
		return respondJSON(http.StatusMethodNotAllowed, ErrorResponse{Message: "Only POST method is allowed"})
	}

	var req TranscodeRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		log.Printf("Invalid request body: %v | raw body: %s", err, event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
	}

	if strings.TrimSpace(req.AssetID) == "" || strings.TrimSpace(req.VideoType) == "" || strings.TrimSpace(req.Format) == "" {
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "assetId, videoType, and format are required"})
	}

	if req.Format != "hls" && req.Format != "dash" {
		return respondJSON(http.StatusBadRequest, ErrorResponse{Message: "format must be either 'hls' or 'dash'"})
	}

	queueURL := os.Getenv("TRANSCODER_QUEUE_URL")
	if queueURL == "" {
		log.Println("Missing TRANSCODER_QUEUE_URL env variable")
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server configuration error: missing queue URL"})
	}

	region := os.Getenv("AWS_REGION")
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
		log.Printf("Failed to create AWS session: %v", err)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to initialize AWS session"})
	}

	svc := sqs.New(sess)

	jobType := "transcode-" + req.Format

	var outputBucketName string
	switch req.Format {
	case "hls":
		outputBucketName = "hls-storage"
	case "dash":
		outputBucketName = "dash-storage"
	}

	var outputFilename string
	if req.SourceFileName != "" {
		nameWithoutExt := strings.TrimSuffix(req.SourceFileName, filepath.Ext(req.SourceFileName))
		if nameWithoutExt != "" {
			switch req.Format {
			case "hls":
				outputFilename = nameWithoutExt + ".m3u8"
			case "dash":
				outputFilename = nameWithoutExt + ".mpd"
			}
		}
	}

	if outputFilename == "" {
		log.Printf("Warning: SourceFileName is empty, using fallback filename for asset %s", req.AssetID)
		switch req.Format {
		case "hls":
			outputFilename = "playlist.m3u8"
		case "dash":
			outputFilename = "manifest.mpd"
		}
	}

	outputKey := fmt.Sprintf("%s/%s/%s", req.AssetID, req.VideoType, outputFilename)
	outputPath := fmt.Sprintf("s3://%s/%s", outputBucketName, outputKey)

	payload := map[string]interface{}{
		"assetId":        req.AssetID,
		"videoType":      req.VideoType,
		"format":         req.Format,
		"input":          req.Input,
		"output":         outputPath,
		"outputBucket":   outputBucketName,
		"outputKey":      outputKey,
		"outputFileName": outputFilename,
	}

	messageBody, err := json.Marshal(map[string]interface{}{
		"type":    jobType,
		"payload": payload,
	})
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to create job message"})
	}

	log.Printf("Sending SQS message for %s: %s", jobType, string(messageBody))

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	}

	_, err = svc.SendMessageWithContext(ctx, input)
	if err != nil {
		log.Printf("Failed to send SQS message: %v", err)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to trigger transcoding job"})
	}

	log.Printf("Transcoding job triggered successfully: %s for asset %s", jobType, req.AssetID)

	return respondJSON(http.StatusOK, TranscodeResponse{
		Message: "Transcoding job triggered successfully",
		JobType: jobType,
	})
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
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
