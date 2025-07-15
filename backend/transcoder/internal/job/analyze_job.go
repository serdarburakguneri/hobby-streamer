package job

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type AnalyzeRunner struct {
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
	s3Client        *s3.Client
}

func NewAnalyzeRunner() *AnalyzeRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &AnalyzeRunner{
		logger:   logger.WithService("analyze-runner"),
		s3Client: s3Client,
	}
}

func NewAnalyzeRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *AnalyzeRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &AnalyzeRunner{
		logger:          logger.WithService("analyze-runner"),
		analyzeProducer: analyzeProducer,
		s3Client:        s3Client,
	}
}

func (a *AnalyzeRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := a.logger.WithContext(ctx)

	var p messages.AnalyzePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal analyze payload")
		return err
	}

	log.Info("Starting video analysis", "input", p.Input, "asset_id", p.AssetID, "video_type", p.VideoType)

	var localPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localPath, err = a.s3Client.Download(ctx, p.Input)
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

	payload := messages.AnalyzeCompletionPayload{
		AssetID:   assetID,
		VideoType: videoType,
		Success:   success,
	}

	if !success && errorMessage != "" {
		payload.Error = errorMessage
	}

	err := a.analyzeProducer.SendMessage(ctx, messages.MessageTypeAnalyzeCompleted, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze completed message", "asset_id", assetID, "success", success)
	} else {
		log.Info("Analyze completed message sent successfully", "asset_id", assetID, "success", success)
	}
}
