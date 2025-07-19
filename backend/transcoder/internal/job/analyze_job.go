package job

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
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
	metadata           *VideoMetadata
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

	err := apperrors.Retry(ctx, func(ctx context.Context) error {
		return a.executeAnalysis(ctx, p)
	}, apperrors.DefaultRetryConfig())

	if err != nil {
		log.WithError(err).Error("Video analysis failed after retries", "asset_id", p.AssetID, "video_id", p.VideoID)
		if a.completionProducer != nil {
			a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, false, "analysis failed after retries")
		}
		return err
	}

	if a.completionProducer != nil {
		a.sendAnalyzeCompleted(ctx, p.AssetID, p.VideoID, true, "")
	}

	return nil
}

type VideoMetadata struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Duration    float64 `json:"duration"`
	Bitrate     int     `json:"bitrate"`
	Codec       string  `json:"codec"`
	Size        int64   `json:"size"`
	ContentType string  `json:"contentType"`
}

func (a *AnalyzeRunner) executeAnalysis(ctx context.Context, p messages.AnalyzePayload) error {
	var localPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localPath, err = a.s3Client.Download(ctx, p.Input)
		if err != nil {
			a.logger.WithError(err).Error("Failed to download from S3", "input", p.Input, "asset_id", p.AssetID, "video_id", p.VideoID)
			return apperrors.NewExternalError("failed to download input file from S3", err)
		}
		defer os.Remove(localPath)
	} else {
		localPath = p.Input
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			a.logger.Error("Input file does not exist", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID)
			return apperrors.NewNotFoundError("input file not found", err)
		}
	}

	metadata, err := a.extractVideoMetadata(ctx, localPath)
	if err != nil {
		a.logger.WithError(err).Error("Failed to extract video metadata", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID)
		return apperrors.NewInternalError("failed to extract video metadata", err)
	}

	a.metadata = metadata

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localPath, "-f", "null", "-")
	out, err := cmd.CombinedOutput()

	if err != nil {
		a.logger.WithError(err).Error("FFmpeg analysis failed", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID, "ffmpeg_output", string(out))
		return apperrors.NewInternalError("ffmpeg analysis failed", err)
	}

	a.logger.Info("Video analysis completed successfully", "input", localPath, "asset_id", p.AssetID, "video_id", p.VideoID, "output_length", len(out), "metadata", metadata)
	a.logger.Debug("FFmpeg analysis output", "output", string(out))

	return nil
}

func (a *AnalyzeRunner) extractVideoMetadata(ctx context.Context, filePath string) (*VideoMetadata, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var probeResult struct {
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
			Size     string `json:"size"`
		} `json:"format"`
		Streams []struct {
			CodecType string `json:"codec_type"`
			CodecName string `json:"codec_name"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out, &probeResult); err != nil {
		return nil, err
	}

	metadata := &VideoMetadata{}

	for _, stream := range probeResult.Streams {
		if stream.CodecType == "video" {
			metadata.Width = stream.Width
			metadata.Height = stream.Height
			metadata.Codec = stream.CodecName
			break
		}
	}

	if probeResult.Format.Duration != "" {
		if duration, err := strconv.ParseFloat(probeResult.Format.Duration, 64); err == nil {
			metadata.Duration = duration
		}
	}

	if probeResult.Format.BitRate != "" {
		if bitrate, err := strconv.Atoi(probeResult.Format.BitRate); err == nil {
			metadata.Bitrate = bitrate
		}
	}

	if probeResult.Format.Size != "" {
		if size, err := strconv.ParseInt(probeResult.Format.Size, 10, 64); err == nil {
			metadata.Size = size
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		metadata.ContentType = "video/mp4"
	}

	return metadata, nil
}

func (a *AnalyzeRunner) sendAnalyzeCompleted(ctx context.Context, assetID, videoID string, success bool, errorMessage string) {
	log := a.logger.WithContext(ctx)

	payload := messages.AnalyzeCompletionPayload{
		AssetID: assetID,
		VideoID: videoID,
		Success: success,
	}

	if success && a.metadata != nil {
		payload.Width = a.metadata.Width
		payload.Height = a.metadata.Height
		payload.Duration = a.metadata.Duration
		payload.Bitrate = a.metadata.Bitrate
		payload.Codec = a.metadata.Codec
		payload.Size = a.metadata.Size
		payload.ContentType = a.metadata.ContentType
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
