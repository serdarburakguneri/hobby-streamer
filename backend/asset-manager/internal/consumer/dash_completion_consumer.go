package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type DASHCompletionConsumer struct {
	assetService asset.AssetService
	logger       *logger.Logger
}

func NewDASHCompletionConsumer(assetService asset.AssetService) *DASHCompletionConsumer {
	return &DASHCompletionConsumer{
		assetService: assetService,
		logger:       logger.Get().WithService("dash-completion-consumer"),
	}
}

func (d *DASHCompletionConsumer) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := d.logger.WithContext(ctx)

	log.Info("Processing DASH transcoding completion message", "message_type", msgType)

	if payload == nil {
		log.Error("Received nil payload for DASH transcoding completion message")
		return apperrors.NewValidationError("payload cannot be nil", nil)
	}

	assetID, ok := payload["assetId"].(string)
	if !ok || assetID == "" {
		log.Error("Missing or invalid assetId in DASH transcoding completion payload")
		return apperrors.NewValidationError("missing or invalid assetId in payload", nil)
	}

	videoID, ok := payload["videoId"].(string)
	if !ok || videoID == "" {
		log.Error("Missing or invalid videoId in DASH transcoding completion payload")
		return apperrors.NewValidationError("missing or invalid videoId in payload", nil)
	}

	format, ok := payload["format"].(string)
	if !ok || format == "" {
		log.Error("Missing or invalid format in DASH transcoding completion payload")
		return apperrors.NewValidationError("missing or invalid format in payload", nil)
	}

	log.Info("Processing DASH transcoding completion", "asset_id", assetID, "video_id", videoID, "format", format)

	err := d.assetService.HandleTranscodeCompletion(ctx, payload)
	if err != nil {
		log.WithError(err).Error("Failed to handle DASH transcoding completion",
			"asset_id", assetID,
			"video_id", videoID,
			"format", format,
			"message_type", msgType)

		return apperrors.WrapWithContext(err, "failed to process DASH transcoding completion")
	}

	log.Info("Successfully processed DASH transcoding completion", "asset_id", assetID, "video_id", videoID, "format", format)
	return nil
}
