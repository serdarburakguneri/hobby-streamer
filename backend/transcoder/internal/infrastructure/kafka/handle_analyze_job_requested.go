package kafka

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

type JobAnalyzeRequestedEvent struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Input   string `json:"input"`
}

func (c *TranscoderEventConsumer) HandleAnalyzeJobRequested(ctx context.Context, event *events.Event) error {
	c.logger.Info("Analyze job requested event received", "event_id", event.ID, "source", event.Source)

	var e JobAnalyzeRequestedEvent
	if err := c.unmarshalEventData(event, &e); err != nil {
		c.logger.WithError(err).Error("Failed to unmarshal analyze job event")
		return err
	}

	payload := messages.JobPayload{
		JobType: "analyze",
		AssetID: e.AssetID,
		VideoID: e.VideoID,
		Input:   e.Input,
	}

	if err := c.jobService.ProcessJob(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to process analyze job", "asset_id", e.AssetID, "video_id", e.VideoID)
		return err
	}

	c.logger.Info("Analyze job processed successfully", "asset_id", e.AssetID, "video_id", e.VideoID)
	return nil
}
