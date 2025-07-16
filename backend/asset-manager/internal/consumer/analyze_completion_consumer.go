package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AnalyzeCompletionConsumer struct {
	assetService asset.AssetService
	logger       *logger.Logger
}

func NewAnalyzeCompletionConsumer(assetService asset.AssetService) *AnalyzeCompletionConsumer {
	return &AnalyzeCompletionConsumer{
		assetService: assetService,
		logger:       logger.Get().WithService("analyze-completion-consumer"),
	}
}

func (a *AnalyzeCompletionConsumer) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := a.logger.WithContext(ctx)

	log.Info("Processing analyze completion message", "message_type", msgType)

	if payload == nil {
		log.Error("Received nil payload for analyze completion message")
		return apperrors.NewValidationError("payload cannot be nil", nil)
	}

	assetID, ok := payload["assetId"].(string)
	if !ok || assetID == "" {
		log.Error("Missing or invalid assetId in analyze completion payload")
		return apperrors.NewValidationError("missing or invalid assetId in payload", nil)
	}

	videoID, ok := payload["videoId"].(string)
	if !ok || videoID == "" {
		log.Error("Missing or invalid videoId in analyze completion payload")
		return apperrors.NewValidationError("missing or invalid videoId in payload", nil)
	}

	log.Info("Processing analyze completion", "asset_id", assetID, "video_id", videoID)

	err := a.assetService.HandleAnalyzeCompletion(ctx, payload)
	if err != nil {
		log.WithError(err).Error("Failed to handle analyze completion",
			"asset_id", assetID,
			"video_id", videoID,
			"message_type", msgType)

		errorType := apperrors.GetErrorType(err)
		switch errorType {
		case apperrors.ErrorTypeValidation:
			return apperrors.NewValidationError("invalid analyze completion payload", err)
		case apperrors.ErrorTypeNotFound:
			return apperrors.NewNotFoundError("asset or video not found for analyze completion", err)
		case apperrors.ErrorTypeTransient:
			return apperrors.NewTransientError("temporary failure processing analyze completion", err)
		default:
			return apperrors.NewInternalError("failed to process analyze completion", err)
		}
	}

	log.Info("Successfully processed analyze completion", "asset_id", assetID, "video_id", videoID)
	return nil
}
