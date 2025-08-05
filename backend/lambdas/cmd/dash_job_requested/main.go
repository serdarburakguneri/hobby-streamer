package main

import (
	"context"
	"encoding/json"
	"fmt"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	pkgevents "github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Input struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Input   string `json:"input"`
}

type DASHJobRequestedEvent struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Input   string `json:"input"`
}

func corsHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "http://localhost:8081",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
	}
}

func NewDASHJobRequestedEvent(assetID, videoID, input string) *DASHJobRequestedEvent {
	return &DASHJobRequestedEvent{
		AssetID: assetID,
		VideoID: videoID,
		Input:   input,
	}
}

func (e *DASHJobRequestedEvent) ToCloudEvent() *pkgevents.Event {
	event := pkgevents.NewEvent("dash.job.requested", e)
	event.SetSource("cms-lambda")
	return event
}

func handleHTTPRequest(ctx context.Context, request awsevents.APIGatewayProxyRequest) (awsevents.APIGatewayProxyResponse, error) {
	logger.Init(logger.GetLogLevel("INFO"), "json")
	log := logger.WithService("dash-job-requested-lambda")

	bootstrap := "kafka:29092"

	producer, err := pkgevents.NewProducer(ctx, &pkgevents.ProducerConfig{
		BootstrapServers: []string{bootstrap},
		Source:           "cms-lambda",
		MaxMessageBytes:  1000000,
	})
	if err != nil {
		log.WithError(err).Error("Failed to create Kafka producer")
		return awsevents.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    corsHeaders(),
			Body:       fmt.Sprintf(`{"error": "Failed to create Kafka producer: %v"}`, err),
		}, nil
	}

	var input Input
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		log.WithError(err).Error("Failed to unmarshal request body")
		return awsevents.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    corsHeaders(),
			Body:       fmt.Sprintf(`{"error": "Invalid request body: %v"}`, err),
		}, nil
	}

	if input.AssetID == "" || input.VideoID == "" || input.Input == "" {
		log.Error("Missing required fields", "asset_id", input.AssetID, "video_id", input.VideoID, "input", input.Input)
		return awsevents.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    corsHeaders(),
			Body:       `{"error": "Missing required fields: assetId, videoId, input"}`,
		}, nil
	}

	event := NewDASHJobRequestedEvent(input.AssetID, input.VideoID, input.Input)

	if err := producer.SendEvent(ctx, "dash.job.requested", event.ToCloudEvent()); err != nil {
		log.WithError(err).Error("Failed to send event", "asset_id", input.AssetID, "video_id", input.VideoID)
		return awsevents.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    corsHeaders(),
			Body:       fmt.Sprintf(`{"error": "Failed to send event: %v"}`, err),
		}, nil
	}

	log.Info("DASH job requested event sent successfully", "asset_id", input.AssetID, "video_id", input.VideoID)

	return awsevents.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    corsHeaders(),
		Body:       fmt.Sprintf(`{"message": "DASH job requested successfully", "assetId": "%s", "videoId": "%s"}`, input.AssetID, input.VideoID),
	}, nil
}

func main() {
	lambda.Start(handleHTTPRequest)
}
