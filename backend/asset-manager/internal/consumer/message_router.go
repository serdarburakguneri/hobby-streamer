package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
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

func (r *MessageRouter) HandleAnalyzeMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := r.logger.WithContext(ctx)
	log.Info("Handling analyze message", "message_type", msgType)
	return r.analyzeConsumer.HandleMessage(ctx, msgType, payload)
}

func (r *MessageRouter) HandleHLSMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := r.logger.WithContext(ctx)
	log.Info("Handling HLS message", "message_type", msgType)
	return r.hlsConsumer.HandleMessage(ctx, msgType, payload)
}

func (r *MessageRouter) HandleDASHMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := r.logger.WithContext(ctx)
	log.Info("Handling DASH message", "message_type", msgType)
	return r.dashConsumer.HandleMessage(ctx, msgType, payload)
}

func (r *MessageRouter) HandleMessage(ctx context.Context, msgType string, payload map[string]interface{}) error {
	log := r.logger.WithContext(ctx)

	log.Info("Routing message", "message_type", msgType)

	if msgType == "" {
		log.Error("Received message with empty message type")
		return apperrors.NewValidationError("message type cannot be empty", nil)
	}

	var err error
	switch msgType {
	case messages.MessageTypeAnalyzeCompleted:
		log.Info("Routing to analyze completion consumer", "message_type", msgType)
		err = r.analyzeConsumer.HandleMessage(ctx, msgType, payload)
	case messages.MessageTypeTranscodeHLSCompleted:
		log.Info("Routing to HLS completion consumer", "message_type", msgType)
		err = r.hlsConsumer.HandleMessage(ctx, msgType, payload)
	case messages.MessageTypeTranscodeDASHCompleted:
		log.Info("Routing to DASH completion consumer", "message_type", msgType)
		err = r.dashConsumer.HandleMessage(ctx, msgType, payload)
	default:
		log.Warn("Unknown message type, ignoring", "message_type", msgType)
		return nil
	}

	if err != nil {
		log.WithError(err).Error("Failed to handle message", "message_type", msgType)

		errorType := apperrors.GetErrorType(err)
		switch errorType {
		case apperrors.ErrorTypeValidation:
			return apperrors.NewValidationError("invalid message payload", err)
		case apperrors.ErrorTypeNotFound:
			return apperrors.NewNotFoundError("resource not found for message processing", err)
		case apperrors.ErrorTypeTransient:
			return apperrors.NewTransientError("temporary failure processing message", err)
		default:
			return apperrors.NewInternalError("failed to process message", err)
		}
	}

	log.Info("Successfully routed and processed message", "message_type", msgType)
	return nil
}
