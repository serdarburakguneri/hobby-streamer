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
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
	outputBucket    string
	outputKey       string
	outputFileName  string
	s3Client        *s3.Client
}

func NewTranscodeDASHRunner() *TranscodeDASHRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeDASHRunner{
		logger:   logger.WithService("dash-runner"),
		s3Client: s3Client,
	}
}

func NewTranscodeDASHRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *TranscodeDASHRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeDASHRunner{
		logger:          logger.WithService("dash-runner"),
		analyzeProducer: analyzeProducer,
		s3Client:        s3Client,
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
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "invalid payload format")
		}
		return apperrors.NewValidationError("invalid payload format", err)
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

	var localInputPath string
	var localOutputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = d.s3Client.Download(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download input from S3", "input", p.Input, "asset_id", p.AssetID, "video_id", p.VideoID)
			if d.analyzeProducer != nil {
				d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "failed to download input file")
			}
			return apperrors.NewExternalError("failed to download input file from S3", err)
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
		if _, err := os.Stat(localInputPath); os.IsNotExist(err) {
			log.Error("Input file does not exist", "input", localInputPath, "asset_id", p.AssetID, "video_id", p.VideoID)
			if d.analyzeProducer != nil {
				d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "input file not found")
			}
			return apperrors.NewNotFoundError("input file not found", err)
		}
	}

	tempDir := os.TempDir()
	dashDir := filepath.Join(tempDir, "dash_output")
	if err := os.MkdirAll(dashDir, 0755); err != nil {
		log.WithError(err).Error("Failed to create DASH output directory", "dir", dashDir, "asset_id", p.AssetID, "video_id", p.VideoID)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "failed to create output directory")
		}
		return apperrors.NewInternalError("failed to create output directory", err)
	}
	defer os.RemoveAll(dashDir)

	localOutputPath = filepath.Join(dashDir, d.outputFileName)

	log.Info("Running FFmpeg DASH transcoding", "input", localInputPath, "output", localOutputPath, "asset_id", p.AssetID, "video_id", p.VideoID)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localInputPath, "-c:v", "libx264", "-c:a", "aac", "-f", "dash", localOutputPath)
	cmd.Dir = dashDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg DASH transcoding failed", "input", localInputPath, "output", localOutputPath, "asset_id", p.AssetID, "video_id", p.VideoID, "ffmpeg_output", string(out))
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "ffmpeg transcoding failed")
		}
		return apperrors.NewInternalError("ffmpeg transcoding failed", err)
	}

	log.Info("DASH transcoding completed successfully", "input", localInputPath, "output", localOutputPath, "asset_id", p.AssetID, "video_id", p.VideoID, "output_length", len(out))
	log.Debug("FFmpeg DASH output", "output", string(out))

	entries, err := os.ReadDir(dashDir)
	if err != nil {
		log.WithError(err).Error("Failed to read DASH output directory", "dir", dashDir, "asset_id", p.AssetID, "video_id", p.VideoID)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "failed to read output directory")
		}
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
				log.WithError(upErr).Error("Failed to upload file to S3", "bucket", d.outputBucket, "key", s3Key, "asset_id", p.AssetID, "video_id", p.VideoID)
				uploadErrors++
				if d.analyzeProducer != nil {
					d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "failed to upload transcoded files")
				}
				return apperrors.NewExternalError("failed to upload transcoded files to S3", upErr)
			}
		}
	}

	if uploadErrors > 0 {
		log.Error("Some files failed to upload", "upload_errors", uploadErrors, "asset_id", p.AssetID, "video_id", p.VideoID)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, "partial upload failure")
		}
		return apperrors.NewExternalError("partial upload failure", nil)
	}

	if d.analyzeProducer != nil {
		d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, true, "")
	}

	log.Info("DASH transcoding job completed successfully", "asset_id", p.AssetID, "video_id", p.VideoID, "format", p.Format)
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

	var messageType string
	switch format {
	case "hls":
		messageType = messages.MessageTypeTranscodeHLSCompleted
	case "dash":
		messageType = messages.MessageTypeTranscodeDASHCompleted
	default:
		log.Error("Unknown format for transcode completion", "format", format, "asset_id", assetID, "video_id", videoID)
		return
	}

	err := d.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "video_id", videoID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "video_id", videoID, "format", format, "success", success)
	}
}
