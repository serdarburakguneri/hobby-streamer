package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type StatusConsumer struct {
	assetService asset.AssetService
	logger       *logger.Logger
}

func NewStatusConsumer(assetService asset.AssetService) *StatusConsumer {
	return &StatusConsumer{
		assetService: assetService,
		logger:       logger.Get().WithService("status-consumer"),
	}
}

func (s *StatusConsumer) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	s.logger.Info("Processing status message", "message_type", msgType)
	return s.assetService.HandleStatusUpdateMessage(ctx, msgType, payload)
}
