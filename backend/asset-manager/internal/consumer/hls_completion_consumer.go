package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
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
	h.logger.Info("Processing HLS transcoding completion message", "message_type", msgType)
	return h.assetService.HandleTranscodeCompletion(ctx, payload, asset.VideoVariantHLS)
}
