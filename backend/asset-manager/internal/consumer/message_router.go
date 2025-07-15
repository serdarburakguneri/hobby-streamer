package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

type MessageRouter struct {
	analyzeConsumer *AnalyzeCompletionConsumer
	hlsConsumer     *HLSCompletionConsumer
	dashConsumer    *DASHCompletionConsumer
	logger          *logger.Logger
}

func NewMessageRouter(assetService asset.AssetService) *MessageRouter {
	return &MessageRouter{
		analyzeConsumer: NewAnalyzeCompletionConsumer(assetService),
		hlsConsumer:     NewHLSCompletionConsumer(assetService),
		dashConsumer:    NewDASHCompletionConsumer(assetService),
		logger:          logger.Get().WithService("message-router"),
	}
}

func (r *MessageRouter) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	r.logger.Info("Routing message", "message_type", msgType)

	switch msgType {
	case messages.MessageTypeAnalyzeCompleted:
		return r.analyzeConsumer.HandleMessage(ctx, msgType, payload)
	case messages.MessageTypeTranscodeHLSCompleted:
		return r.hlsConsumer.HandleMessage(ctx, msgType, payload)
	case messages.MessageTypeTranscodeDASHCompleted:
		return r.dashConsumer.HandleMessage(ctx, msgType, payload)
	default:
		r.logger.Info("Unknown message type, ignoring", "message_type", msgType)
		return nil
	}
}
