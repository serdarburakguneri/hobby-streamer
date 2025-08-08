package kafka

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

type HLSJobRequestedEvent struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Input   string `json:"input"`
	JobID   string `json:"jobId,omitempty"`
}

func (c *TranscoderEventConsumer) HandleHLSJobRequested(ctx context.Context, event *events.Event) error {
	c.logger.Info("HLS job requested event received", "event_id", event.ID, "source", event.Source)

	var e HLSJobRequestedEvent
	if err := c.unmarshalEventData(event, &e); err != nil {
		c.logger.WithError(err).Error("Failed to unmarshal HLS job event")
		return err
	}

	payload := messages.JobPayload{
		JobID:   e.JobID,
		JobType: "transcode",
		AssetID: e.AssetID,
		VideoID: e.VideoID,
		Input:   e.Input,
		Format:  "hls",
		Quality: "main",
	}

	if err := c.jobService.ProcessJob(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to process HLS job", "asset_id", e.AssetID, "video_id", e.VideoID)
		return err
	}

	c.logger.Info("HLS job processed successfully", "asset_id", e.AssetID, "video_id", e.VideoID)
	return nil
}
