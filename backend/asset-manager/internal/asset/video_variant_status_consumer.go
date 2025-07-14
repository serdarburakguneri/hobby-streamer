package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AnalyzeCompletionConsumer struct {
	assetService AssetService
	logger       *logger.Logger
}

func NewAnalyzeCompletionConsumer(assetService AssetService) *AnalyzeCompletionConsumer {
	return &AnalyzeCompletionConsumer{
		assetService: assetService,
		logger:       logger.Get().WithService("analyze-completion-consumer"),
	}
}

func (a *AnalyzeCompletionConsumer) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	a.logger.Info("Processing analyze completion message", "message_type", msgType)
	return a.assetService.HandleAnalyzeCompletionMessage(ctx, msgType, payload)
}
