package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

type TranscodeRequest struct {
	AssetID        string `json:"assetId"`
	VideoID        string `json:"videoId"`
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
	Type    string `json:"type,omitempty"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	trackingID := ""
	if trackingHeader, exists := event.Headers["X-Tracking-ID"]; exists {
		trackingID = trackingHeader
	} else if trackingHeader, exists := event.Headers["x-tracking-id"]; exists {
		trackingID = trackingHeader
	}

	log := logger.WithService("trigger-transcode-job")
	if trackingID != "" {
		log = log.WithTrackingID(trackingID)
	}
	log = log.WithContext(ctx)

	if event.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost:8081",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With",
				"Access-Control-Allow-Credentials": "true",
			},
		}, nil
	}

	if event.HTTPMethod != "POST" {
		log.Error("Invalid HTTP method", "method", event.HTTPMethod)
		return respondJSON(http.StatusMethodNotAllowed, ErrorResponse{
			Message: "Only POST method is allowed",
			Type:    "validation",
		})
	}

	if event.Body == "" {
		log.Error("Empty request body")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Request body is required",
			Type:    "validation",
		})
	}

	var req TranscodeRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		log.WithError(err).Error("Invalid request body", "raw_body", event.Body)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request body format",
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

	if strings.TrimSpace(req.VideoID) == "" {
		log.Error("Missing videoId in request")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "videoId is required",
			Type:    "validation",
		})
	}

	if strings.TrimSpace(req.Format) == "" {
		log.Error("Missing format in request")
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "format is required",
			Type:    "validation",
		})
	}

	if req.Format != "hls" && req.Format != "dash" {
		log.Error("Invalid format", "format", req.Format, "asset_id", req.AssetID, "video_id", req.VideoID)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "format must be either 'hls' or 'dash'",
			Type:    "validation",
		})
	}

	if strings.TrimSpace(req.Input) == "" {
		log.Error("Missing input in request", "asset_id", req.AssetID, "video_id", req.VideoID)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "input is required",
			Type:    "validation",
		})
	}

	queueURL := os.Getenv("TRANSCODER_QUEUE_URL")
	if queueURL == "" {
		log.Error("Missing TRANSCODER_QUEUE_URL env variable")
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Server configuration error: missing queue URL",
			Type:    "internal",
		})
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
		log.WithError(err).Error("Failed to create AWS session", "asset_id", req.AssetID, "video_id", req.VideoID)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to initialize AWS session",
			Type:    "internal",
		})
	}

	svc := sqs.New(sess)

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
		log.Warn("SourceFileName is empty, using fallback filename", "asset_id", req.AssetID, "video_id", req.VideoID)
		switch req.Format {
		case "hls":
			outputFilename = "playlist.m3u8"
		case "dash":
			outputFilename = "manifest.mpd"
		}
	}

	outputKey := fmt.Sprintf("%s/%s/%s", req.AssetID, req.VideoID, outputFilename)

	payload := messages.TranscodePayload{
		AssetID:        req.AssetID,
		VideoID:        req.VideoID,
		Format:         req.Format,
		Input:          req.Input,
		OutputBucket:   outputBucketName,
		OutputKey:      outputKey,
		OutputFileName: outputFilename,
	}

	var messageType string
	switch req.Format {
	case "hls":
		messageType = messages.MessageTypeTranscodeHLS
	case "dash":
		messageType = messages.MessageTypeTranscodeDASH
	default:
		log.Error("Unknown format", "format", req.Format, "asset_id", req.AssetID, "video_id", req.VideoID)
		return respondJSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid format",
			Type:    "validation",
		})
	}

	messageBody, err := json.Marshal(map[string]interface{}{
		"type":    messageType,
		"payload": payload,
	})
	if err != nil {
		log.WithError(err).Error("Failed to marshal message", "asset_id", req.AssetID, "video_id", req.VideoID)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to create job message",
			Type:    "internal",
		})
	}

	log.Info("Sending SQS message", "message_type", messageType, "asset_id", req.AssetID, "video_id", req.VideoID, "output_bucket", outputBucketName, "output_key", outputKey)

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	}

	_, err = svc.SendMessageWithContext(ctx, input)
	if err != nil {
		log.WithError(err).Error("Failed to send SQS message", "queue_url", queueURL, "asset_id", req.AssetID, "video_id", req.VideoID)
		return respondJSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to trigger transcoding job",
			Type:    "external",
		})
	}

	log.Info("Transcoding job triggered successfully", "message_type", messageType, "asset_id", req.AssetID, "video_id", req.VideoID)

	response, err := respondJSON(http.StatusOK, TranscodeResponse{
		Message: "Transcoding job triggered successfully",
		JobType: messageType,
	})
	if err != nil {
		return response, err
	}

	if trackingID != "" {
		response.Headers["X-Tracking-ID"] = trackingID
	}

	return response, nil
}

func respondJSON(status int, payload interface{}) (events.APIGatewayProxyResponse, error) {
	log := logger.WithService("trigger-transcode-job")

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
