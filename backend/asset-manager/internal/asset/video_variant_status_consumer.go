package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type VideoVariantStatusConsumer struct {
	assetService AssetService
	logger       *logger.Logger
}

func NewVideoVariantStatusConsumer(assetService AssetService) *VideoVariantStatusConsumer {
	return &VideoVariantStatusConsumer{
		assetService: assetService,
		logger:       logger.Get().WithService("video-variant-status-consumer"),
	}
}

func (v *VideoVariantStatusConsumer) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	v.logger.Info("Processing video variant status message", "message_type", msgType)
	return v.assetService.HandleStatusUpdateMessage(ctx, msgType, payload)
}
