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

type DASHPayload struct {
	Input          string `json:"input"`
	AssetID        string `json:"assetId"`
	VideoType      string `json:"videoType"`
	Format         string `json:"format"`
	OutputBucket   string `json:"outputBucket"`
	OutputKey      string `json:"outputKey"`
	OutputFileName string `json:"outputFileName"`
}

type TranscodeDASHRunner struct {
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
	outputBucket    string
	outputKey       string
	outputFileName  string
}

func NewTranscodeDASHRunner() *TranscodeDASHRunner {
	return &TranscodeDASHRunner{
		logger: logger.WithService("dash-runner"),
	}
}

func NewTranscodeDASHRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *TranscodeDASHRunner {
	return &TranscodeDASHRunner{
		logger:          logger.WithService("dash-runner"),
		analyzeProducer: analyzeProducer,
	}
}

func (d *TranscodeDASHRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := d.logger.WithContext(ctx)

	var p DASHPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal DASH payload")
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	d.outputBucket = p.OutputBucket
	d.outputKey = p.OutputKey
	d.outputFileName = p.OutputFileName

	var localInputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = d.downloadFromS3(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download input from S3", "input", p.Input)
			if d.analyzeProducer != nil {
				d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
			}
			return err
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
	}

	tempDir := os.TempDir()
	outputDir := filepath.Join(tempDir, "dash_output")
	os.MkdirAll(outputDir, 0755)

	manifestPath := filepath.Join(outputDir, d.outputFileName)
	segmentPattern := filepath.Join(outputDir, "segment_%03d.m4s")

	log.Info("Running FFmpeg DASH transcoding", "input", localInputPath, "manifest", manifestPath, "segment_pattern", segmentPattern)
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", localInputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "dash",
		"-init_seg_name", "init_$RepresentationID$.m4s",
		"-media_seg_name", "segment_$RepresentationID$_$Number%03d$.m4s",
		manifestPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg DASH transcoding failed", "input", localInputPath, "manifest", manifestPath, "ffmpeg_output", string(out))
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	log.Info("DASH transcoding completed successfully", "input", localInputPath, "manifest", manifestPath, "output_length", len(out))
	log.Debug("FFmpeg DASH output", "output", string(out))

	err = d.uploadToS3(ctx, manifestPath, d.outputBucket, d.outputKey)
	if err != nil {
		log.WithError(err).Error("Failed to upload DASH manifest to S3", "bucket", d.outputBucket, "key", d.outputKey)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	segmentFiles, err := filepath.Glob(filepath.Join(outputDir, "*.m4s"))
	if err != nil {
		log.WithError(err).Error("Failed to find DASH segment files")
	} else {
		for _, segmentFile := range segmentFiles {
			segmentName := filepath.Base(segmentFile)
			segmentKey := strings.TrimSuffix(d.outputKey, ".mpd") + "/" + segmentName

			err = d.uploadToS3(ctx, segmentFile, d.outputBucket, segmentKey)
			if err != nil {
				log.WithError(err).Error("Failed to upload DASH segment to S3", "bucket", d.outputBucket, "key", segmentKey, "file", segmentFile)
			} else {
				log.Info("Successfully uploaded DASH segment", "bucket", d.outputBucket, "key", segmentKey)
			}
		}
	}

	// Clean up
	defer os.RemoveAll(outputDir)

	if d.analyzeProducer != nil {
		d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, true, "")
	}

	return nil
}

func (d *TranscodeDASHRunner) downloadFromS3(ctx context.Context, s3URL string) (string, error) {
	log := d.logger.WithContext(ctx)

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

func (d *TranscodeDASHRunner) uploadToS3(ctx context.Context, localPath, bucket, key string) error {
	log := d.logger.WithContext(ctx)

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

func (d *TranscodeDASHRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoType, format string, success bool, errorMessage string) {
	log := d.logger.WithContext(ctx)

	payload := map[string]interface{}{
		"assetId":   assetID,
		"videoType": videoType,
		"format":    format,
		"success":   success,
	}

	if success {
		payload["bucket"] = d.outputBucket
		payload["key"] = d.outputKey
		payload["fileName"] = d.outputFileName
		payload["url"] = fmt.Sprintf("s3://%s/%s", d.outputBucket, d.outputKey)
	}

	if !success && errorMessage != "" {
		payload["error"] = errorMessage
	}

	messageType := "transcode-" + format + "-completed"
	err := d.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "format", format, "success", success)
	}
}
