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

type AnalyzePayload struct {
	Input     string `json:"input"`
	AssetID   string `json:"assetId"`
	VideoType string `json:"videoType"`
}

type AnalyzeRunner struct {
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
}

func NewAnalyzeRunner() *AnalyzeRunner {
	return &AnalyzeRunner{
		logger: logger.WithService("analyze-runner"),
	}
}

func NewAnalyzeRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *AnalyzeRunner {
	return &AnalyzeRunner{
		logger:          logger.WithService("analyze-runner"),
		analyzeProducer: analyzeProducer,
	}
}

func (a *AnalyzeRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := a.logger.WithContext(ctx)

	var p AnalyzePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal analyze payload")
		return err
	}

	log.Info("Starting video analysis", "input", p.Input, "asset_id", p.AssetID, "video_type", p.VideoType)

	var localPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localPath, err = a.downloadFromS3(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download from S3", "input", p.Input)
			if a.analyzeProducer != nil {
				a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoType, false, err.Error())
			}
			return err
		}
		defer os.Remove(localPath)
	} else {
		localPath = p.Input
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localPath, "-f", "null", "-")
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg analysis failed", "input", localPath, "output", string(out))
		if a.analyzeProducer != nil {
			a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoType, false, err.Error())
		}
		return err
	}

	log.Info("Video analysis completed successfully", "input", localPath, "output_length", len(out))
	log.Debug("FFmpeg analysis output", "output", string(out))

	if a.analyzeProducer != nil {
		a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoType, true, "")
	}

	return nil
}

func (a *AnalyzeRunner) sendAnalyzeCompleted(ctx context.Context, assetID, videoType string, success bool, errorMessage string) {
	log := a.logger.WithContext(ctx)

	payload := map[string]interface{}{
		"assetId":   assetID,
		"videoType": videoType,
		"success":   success,
	}

	if !success && errorMessage != "" {
		payload["error"] = errorMessage
	}

	err := a.analyzeProducer.SendMessage(ctx, "analyze-completed", payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze completed message", "asset_id", assetID, "success", success)
	} else {
		log.Info("Analyze completed message sent successfully", "asset_id", assetID, "success", success)
	}
}

func (a *AnalyzeRunner) downloadFromS3(ctx context.Context, s3URL string) (string, error) {
	log := a.logger.WithContext(ctx)

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
