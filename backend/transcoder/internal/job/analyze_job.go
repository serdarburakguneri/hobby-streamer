package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AnalyzePayload struct {
	Input string `json:"input"`
}

type AnalyzeRunner struct {
	logger *logger.Logger
}

func NewAnalyzeRunner() *AnalyzeRunner {
	return &AnalyzeRunner{
		logger: logger.WithService("analyze-runner"),
	}
}

func (a *AnalyzeRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := a.logger.WithContext(ctx)

	var p AnalyzePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal analyze payload")
		return err
	}

	log.Info("Starting video analysis", "input", p.Input)

	var localPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localPath, err = a.downloadFromS3(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download from S3", "input", p.Input)
			return err
		}
		defer os.Remove(localPath)
	} else {
		localPath = p.Input
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg analysis failed", "input", localPath, "output", string(out))
		return err
	}

	log.Info("Video analysis completed successfully", "input", localPath, "output_length", len(out))
	log.Debug("FFmpeg analysis output", "output", string(out))
	return nil
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

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localstack:4566"),
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

	_, err = client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		os.Remove(localPath)
		return "", fmt.Errorf("failed to get object from S3: %w", err)
	}

	log.Info("Successfully downloaded from S3", "local_path", localPath)
	return localPath, nil
}
