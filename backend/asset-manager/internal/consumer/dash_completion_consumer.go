package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
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
	d.logger.Info("Processing DASH transcoding completion message", "message_type", msgType)
	return d.assetService.HandleTranscodeCompletion(ctx, payload, asset.VideoVariantDASH)
}
