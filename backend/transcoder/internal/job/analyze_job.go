package job

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type AnalyzeRunner struct {
	logger             *logger.Logger
	completionProducer *sqs.Producer
	s3Client           *s3.Client
}

func NewAnalyzeRunner() *AnalyzeRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &AnalyzeRunner{
		logger:   logger.WithService("analyze-runner"),
		s3Client: s3Client,
	}
}

func NewAnalyzeRunnerWithCompletionProducer(completionProducer *sqs.Producer) *AnalyzeRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &AnalyzeRunner{
		logger:             logger.WithService("analyze-runner"),
		completionProducer: completionProducer,
		s3Client:           s3Client,
	}
}

func (a *AnalyzeRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := a.logger.WithContext(ctx)

	if payload == nil {
		log.Error("Received nil payload for analyze job")
		return apperrors.NewValidationError("payload cannot be nil", nil)
	}

	var p messages.AnalyzePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal analyze payload")
		if a.completionProducer != nil {
			a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, false, "invalid payload format")
		}
		return apperrors.NewValidationError("invalid payload format", err)
	}

	if p.AssetID == "" {
		log.Error("Missing assetId in analyze payload")
		return apperrors.NewValidationError("assetId is required", nil)
	}

	if p.VideoID == "" {
		log.Error("Missing videoId in analyze payload")
		return apperrors.NewValidationError("videoId is required", nil)
	}

	if p.Input == "" {
		log.Error("Missing input in analyze payload")
		return apperrors.NewValidationError("input is required", nil)
	}

	log.Info("Starting video analysis", "input", p.Input, "asset_id", p.AssetID, "video_id", p.VideoID)

	var localPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localPath, err = a.s3Client.Download(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download from S3", "input", p.Input, "asset_id", p.AssetID, "video_id", p.VideoID)
			if a.completionProducer != nil {
				a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, false, "failed to download input file")
			}
			return apperrors.NewExternalError("failed to download input file from S3", err)
		}
		defer os.Remove(localPath)
	} else {
		localPath = p.Input
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			log.Error("Input file does not exist", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID)
			if a.completionProducer != nil {
				a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, false, "input file not found")
			}
			return apperrors.NewNotFoundError("input file not found", err)
		}
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localPath, "-f", "null", "-")
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg analysis failed", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID, "ffmpeg_output", string(out))
		if a.completionProducer != nil {
			a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, false, "ffmpeg analysis failed")
		}
		return apperrors.NewInternalError("ffmpeg analysis failed", err)
	}

	log.Info("Video analysis completed successfully", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID, "output_length", len(out))
	log.Debug("FFmpeg analysis output", "output", string(out))

	if a.completionProducer != nil {
		a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, true, "")
	}

	return nil
}

func (a *AnalyzeRunner) sendAnalyzeCompleted(ctx context.Context, assetID, videoID string, success bool, errorMessage string) {
	log := a.logger.WithContext(ctx)

	payload := messages.AnalyzeCompletionPayload{
		AssetID: assetID,
		VideoID: videoID,
		Success: success,
	}

	if !success && errorMessage != "" {
		payload.Error = errorMessage
	}

	err := a.completionProducer.SendMessage(ctx, messages.MessageTypeAnalyzeCompleted, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze completed message", "asset_id", assetID, "video_id", videoID, "success", success)
	} else {
		log.Info("Analyze completed message sent successfully", "asset_id", assetID, "video_id", videoID, "success", success)
	}
}
