package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type HLSPayload struct {
	Input          string `json:"input"`
	AssetID        string `json:"assetId"`
	VideoType      string `json:"videoType"`
	Format         string `json:"format"`
	OutputBucket   string `json:"outputBucket"`
	OutputKey      string `json:"outputKey"`
	OutputFileName string `json:"outputFileName"`
}

type TranscodeHLSRunner struct {
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
	outputBucket    string
	outputKey       string
	outputFileName  string
}

func NewTranscodeHLSRunner() *TranscodeHLSRunner {
	return &TranscodeHLSRunner{
		logger: logger.WithService("hls-runner"),
	}
}

func NewTranscodeHLSRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *TranscodeHLSRunner {
	return &TranscodeHLSRunner{
		logger:          logger.WithService("hls-runner"),
		analyzeProducer: analyzeProducer,
	}
}

func (h *TranscodeHLSRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := h.logger.WithContext(ctx)

	var p HLSPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal HLS payload")
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	h.outputBucket = p.OutputBucket
	h.outputKey = p.OutputKey
	h.outputFileName = p.OutputFileName

	var localInputPath string
	var localOutputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = h.downloadFromS3(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download input from S3", "input", p.Input)
			if h.analyzeProducer != nil {
				h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
			}
			return err
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
	}

	tempDir := os.TempDir()
	localOutputPath = filepath.Join(tempDir, h.outputFileName)

	log.Info("Running FFmpeg HLS transcoding", "input", localInputPath, "output", localOutputPath)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localInputPath, "-c:v", "libx264", "-c:a", "aac", "-f", "hls", localOutputPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg HLS transcoding failed", "input", localInputPath, "output", localOutputPath, "ffmpeg_output", string(out))
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	log.Info("HLS transcoding completed successfully", "input", localInputPath, "output", localOutputPath, "output_length", len(out))
	log.Debug("FFmpeg HLS output", "output", string(out))

	err = h.uploadToS3(ctx, localOutputPath, h.outputBucket, h.outputKey)
	if err != nil {
		log.WithError(err).Error("Failed to upload output to S3", "bucket", h.outputBucket, "key", h.outputKey)
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}
	defer os.Remove(localOutputPath)

	if h.analyzeProducer != nil {
		h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, true, "")
	}

	return nil
}

func (h *TranscodeHLSRunner) downloadFromS3(ctx context.Context, s3URL string) (string, error) {
	log := h.logger.WithContext(ctx)

	if !strings.HasPrefix(s3URL, "s3://") {
		return "", fmt.Errorf("invalid S3 URL: %s", s3URL)
	}

	parts := strings.SplitN(s3URL[5:], "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid S3 URL format: %s", s3URL)
	}

	bucket := parts[0]
	key := parts[1]

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		awsEndpoint = "http://localstack:4566"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(awsEndpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)

	tempDir := os.TempDir()
	filename := filepath.Base(key)
	if filename == "" {
		filename = fmt.Sprintf("video_%d.mp4", time.Now().Unix())
	}
	localPath := filepath.Join(tempDir, filename)

	log.Info("Downloading from S3", "bucket", bucket, "key", key, "local_path", localPath)

	file, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	result, err := client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		os.Remove(localPath)
		return "", fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		os.Remove(localPath)
		return "", fmt.Errorf("failed to write file to disk: %w", err)
	}

	log.Info("Successfully downloaded from S3", "local_path", localPath)
	return localPath, nil
}

func (h *TranscodeHLSRunner) uploadToS3(ctx context.Context, localPath, bucket, key string) error {
	log := h.logger.WithContext(ctx)

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		awsEndpoint = "http://localstack:4566"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(awsEndpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer file.Close()

	log.Info("Uploading to S3", "bucket", bucket, "key", key, "local_path", localPath)

	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to S3: %w", err)
	}

	log.Info("Successfully uploaded to S3", "bucket", bucket, "key", key)
	return nil
}

func (h *TranscodeHLSRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoType, format string, success bool, errorMessage string) {
	log := h.logger.WithContext(ctx)

	payload := map[string]interface{}{
		"assetId":   assetID,
		"videoType": videoType,
		"format":    format,
		"success":   success,
	}

	if success {
		payload["bucket"] = h.outputBucket
		payload["key"] = h.outputKey
		payload["fileName"] = h.outputFileName
		payload["url"] = fmt.Sprintf("s3://%s/%s", h.outputBucket, h.outputKey)
	}

	if !success && errorMessage != "" {
		payload["error"] = errorMessage
	}

	messageType := "transcode-" + format + "-completed"
	err := h.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "format", format, "success", success)
	}
}
