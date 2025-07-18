package job

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type TranscodeDASHRunner struct {
	logger             *logger.Logger
	completionProducer *sqs.Producer
	outputBucket       string
	outputKey          string
	outputFileName     string
	s3Client           *s3.Client
}

func NewTranscodeDASHRunner() *TranscodeDASHRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeDASHRunner{
		logger:   logger.WithService("dash-runner"),
		s3Client: s3Client,
	}
}

func NewTranscodeDASHRunnerWithCompletionProducer(completionProducer *sqs.Producer) *TranscodeDASHRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeDASHRunner{
		logger:             logger.WithService("dash-runner"),
		completionProducer: completionProducer,
		s3Client:           s3Client,
	}
}

func (d *TranscodeDASHRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := d.logger.WithContext(ctx)

	if payload == nil {
		log.Error("Received nil payload for DASH transcode job")
		return apperrors.NewValidationError("payload cannot be nil", nil)
	}

	var p messages.TranscodePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal DASH payload")
		return apperrors.NewValidationError("invalid payload format", err)
	}

	if p.Format != "dash" {
		log.Error("Invalid format for DASH transcode job", "format", p.Format)
		return apperrors.NewValidationError("format must be 'dash' for DASH transcode job", nil)
	}

	if p.AssetID == "" {
		log.Error("Missing assetId in DASH transcode payload")
		return apperrors.NewValidationError("assetId is required", nil)
	}

	if p.VideoID == "" {
		log.Error("Missing videoId in DASH transcode payload")
		return apperrors.NewValidationError("videoId is required", nil)
	}

	if p.Input == "" {
		log.Error("Missing input in DASH transcode payload")
		return apperrors.NewValidationError("input is required", nil)
	}

	if p.OutputBucket == "" {
		log.Error("Missing outputBucket in DASH transcode payload")
		return apperrors.NewValidationError("outputBucket is required", nil)
	}

	if p.OutputFileName == "" {
		log.Error("Missing outputFileName in DASH transcode payload")
		return apperrors.NewValidationError("outputFileName is required", nil)
	}

	d.outputBucket = p.OutputBucket
	d.outputKey = p.OutputKey
	d.outputFileName = p.OutputFileName

	log.Info("Starting DASH transcoding", "input", p.Input, "asset_id", p.AssetID, "video_id", p.VideoID, "output_bucket", p.OutputBucket, "output_file", p.OutputFileName)

	err := apperrors.Retry(ctx, func(ctx context.Context) error {
		return d.executeTranscoding(ctx, p)
	}, apperrors.DefaultRetryConfig())

	if err != nil {
		log.WithError(err).Error("DASH transcoding failed after retries", "asset_id", p.AssetID, "video_id", p.VideoID, "format", p.Format)
		if d.completionProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "transcoding failed after retries")
		}
		return err
	}

	if d.completionProducer != nil {
		d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, true, "")
	}

	log.Info("DASH transcoding job completed successfully", "asset_id", p.AssetID, "video_id", p.VideoID, "format", p.Format)
	return nil
}

func (d *TranscodeDASHRunner) executeTranscoding(ctx context.Context, p messages.TranscodePayload) error {
	var localInputPath string
	var localOutputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = d.s3Client.Download(ctx, p.Input)
		if err != nil {
			d.logger.WithError(err).Error("Failed to download input from S3", "input", p.Input, "asset_id", p.AssetID, "video_id", p.VideoID)
			return apperrors.NewExternalError("failed to download input file from S3", err)
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
		if _, err := os.Stat(localInputPath); os.IsNotExist(err) {
			d.logger.Error("Input file does not exist", "input", localInputPath, "asset_id", p.AssetID, "video_id", p.VideoID)
			return apperrors.NewNotFoundError("input file not found", err)
		}
	}

	tempDir := os.TempDir()
	dashDir := filepath.Join(tempDir, "dash_output")
	if err := os.MkdirAll(dashDir, 0750); err != nil {
		d.logger.WithError(err).Error("Failed to create DASH output directory", "dir", dashDir, "asset_id", p.AssetID, "video_id", p.VideoID)
		return apperrors.NewInternalError("failed to create output directory", err)
	}
	defer os.RemoveAll(dashDir)

	localOutputPath = filepath.Join(dashDir, d.outputFileName)

	d.logger.Info("Running FFmpeg DASH transcoding", "input", localInputPath, "output", localOutputPath, "asset_id", p.AssetID, "video_id", p.VideoID)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localInputPath, "-c:v", "libx264", "-c:a", "aac", "-f", "dash", localOutputPath)
	cmd.Dir = dashDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		d.logger.WithError(err).Error("FFmpeg DASH transcoding failed", "input", localInputPath, "output", localOutputPath, "asset_id", p.AssetID, "video_id", p.VideoID, "ffmpeg_output", string(out))
		return apperrors.NewInternalError("ffmpeg transcoding failed", err)
	}

	d.logger.Info("DASH transcoding completed successfully", "input", localInputPath, "output", localOutputPath, "asset_id", p.AssetID, "video_id", p.VideoID, "output_length", len(out))
	d.logger.Debug("FFmpeg DASH output", "output", string(out))

	entries, err := os.ReadDir(dashDir)
	if err != nil {
		d.logger.WithError(err).Error("Failed to read DASH output directory", "dir", dashDir, "asset_id", p.AssetID, "video_id", p.VideoID)
		return apperrors.NewInternalError("failed to read output directory", err)
	}

	uploadErrors := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".mpd") || strings.HasSuffix(name, ".m4s") {
			localPath := filepath.Join(dashDir, name)
			s3Key := filepath.Join(p.AssetID, p.VideoID, name)
			if upErr := d.s3Client.Upload(ctx, localPath, d.outputBucket, s3Key); upErr != nil {
				d.logger.WithError(upErr).Error("Failed to upload file to S3", "bucket", d.outputBucket, "key", s3Key, "asset_id", p.AssetID, "video_id", p.VideoID)
				uploadErrors++
				return apperrors.NewExternalError("failed to upload transcoded files to S3", upErr)
			}
		}
	}

	if uploadErrors > 0 {
		d.logger.Error("Some files failed to upload", "upload_errors", uploadErrors, "asset_id", p.AssetID, "video_id", p.VideoID)
		return apperrors.NewExternalError("partial upload failure", nil)
	}

	return nil
}

func (d *TranscodeDASHRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoID, format string, success bool, errorMessage string) {
	log := d.logger.WithContext(ctx)

	payload := messages.TranscodeCompletionPayload{
		AssetID: assetID,
		VideoID: videoID,
		Format:  format,
		Success: success,
	}

	if success {
		payload.Bucket = d.outputBucket
		payload.Key = d.outputKey
		payload.FileName = d.outputFileName
		payload.URL = "s3://" + d.outputBucket + "/" + d.outputKey
	}

	if !success && errorMessage != "" {
		payload.Error = errorMessage
	}

	messageType := messages.MessageTypeTranscodeDASHCompleted

	err := d.completionProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "video_id", videoID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "video_id", videoID, "format", format, "success", success)
	}
}
