package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	pkgevents "github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type RawVideoUploadedEvent struct {
	AssetID         string  `json:"assetId"`
	VideoID         string  `json:"videoId"`
	StorageLocation string  `json:"storageLocation"`
	Filename        string  `json:"filename"`
	Size            int64   `json:"size"`
	ContentType     string  `json:"contentType"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Duration        float64 `json:"duration"`
	Bitrate         int     `json:"bitrate"`
	Codec           string  `json:"codec"`
}

func NewRawVideoUploadedEvent(assetID, videoID, storageLocation, filename string, size int64, contentType string) *RawVideoUploadedEvent {
	return &RawVideoUploadedEvent{
		AssetID:         assetID,
		VideoID:         videoID,
		StorageLocation: storageLocation,
		Filename:        filename,
		Size:            size,
		ContentType:     contentType,
		Width:           0,
		Height:          0,
		Duration:        0,
		Bitrate:         0,
		Codec:           "",
	}
}

func (e *RawVideoUploadedEvent) ToCloudEvent() *pkgevents.Event {
	event := pkgevents.NewEvent("raw-video-uploaded", e)
	event.SetSource("upload-lambda")
	return event
}

func extractAssetAndVideoID(key string) (string, string, error) {
	parts := strings.Split(key, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid S3 key format: expected assetId/source/filename, got %s", key)
	}
	assetID := parts[0]
	videoID := parts[1] + "/" + parts[2] // source/filename as video ID
	return assetID, videoID, nil
}

func handleS3Event(ctx context.Context, s3Event awsevents.S3Event) error {
	logger.Init(logger.GetLogLevel("INFO"), "json")
	log := logger.WithService("raw-video-uploaded-lambda")

	bootstrap := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if bootstrap == "" {
		bootstrap = "kafka:29092"
	}

	producer, err := pkgevents.NewProducer(ctx, &pkgevents.ProducerConfig{
		BootstrapServers: []string{bootstrap},
		Source:           "upload-lambda",
		MaxMessageBytes:  1000000,
	})
	if err != nil {
		log.WithError(err).Error("Failed to create Kafka producer")
		return err
	}

	for _, record := range s3Event.Records {
		s3 := record.S3
		bucket := s3.Bucket.Name
		key := s3.Object.Key
		size := s3.Object.Size

		log.Info("Processing S3 event", "bucket", bucket, "key", key, "size", size)

		assetID, videoID, err := extractAssetAndVideoID(key)
		if err != nil {
			log.WithError(err).Error("Failed to extract asset and video ID from S3 key", "key", key)
			continue
		}

		storageLocation := fmt.Sprintf("s3://%s/%s", bucket, key)
		filename := s3.Object.Key
		contentType := "video/mp4"

		event := NewRawVideoUploadedEvent(
			assetID,
			videoID,
			storageLocation,
			filename,
			size,
			contentType,
		)

		if err := producer.SendEvent(ctx, "raw-video-uploaded", event.ToCloudEvent()); err != nil {
			log.WithError(err).Error("Failed to send event", "asset_id", assetID, "video_id", videoID)
			return err
		}

		log.Info("Raw video uploaded event sent successfully", "asset_id", assetID, "video_id", videoID)
	}

	return nil
}

func main() {
	lambda.Start(handleS3Event)
}
