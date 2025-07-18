package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type HLSCompletionConsumer struct {
	assetService asset.AssetService
	logger       *logger.Logger
}

func NewHLSCompletionConsumer(assetService asset.AssetService) *HLSCompletionConsumer {
	return &HLSCompletionConsumer{
		assetService: assetService,
		logger:       logger.Get().WithService("hls-completion-consumer"),
	}
}

func (h *HLSCompletionConsumer) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := h.logger.WithContext(ctx)

	log.Info("Processing HLS transcoding completion message", "message_type", msgType)

	if payload == nil {
		log.Error("Received nil payload for HLS transcoding completion message")
		return apperrors.NewValidationError("payload cannot be nil", nil)
	}

	assetID, ok := payload["assetId"].(string)
	if !ok || assetID == "" {
		log.Error("Missing or invalid assetId in HLS transcoding completion payload")
		return apperrors.NewValidationError("missing or invalid assetId in payload", nil)
	}

	videoID, ok := payload["videoId"].(string)
	if !ok || videoID == "" {
		log.Error("Missing or invalid videoId in HLS transcoding completion payload")
		return apperrors.NewValidationError("missing or invalid videoId in payload", nil)
	}

	format, ok := payload["format"].(string)
	if !ok || format == "" {
		log.Error("Missing or invalid format in HLS transcoding completion payload")
		return apperrors.NewValidationError("missing or invalid format in payload", nil)
	}

	log.Info("Processing HLS transcoding completion", "asset_id", assetID, "video_id", videoID, "format", format)

	err := h.assetService.HandleTranscodeCompletion(ctx, payload)
	if err != nil {
		log.WithError(err).Error("Failed to handle HLS transcoding completion",
			"asset_id", assetID,
			"video_id", videoID,
			"format", format,
			"message_type", msgType)

		return apperrors.WrapWithContext(err, "failed to process HLS transcoding completion")
	}

	log.Info("Successfully processed HLS transcoding completion", "asset_id", assetID, "video_id", videoID, "format", format)
	return nil
}
