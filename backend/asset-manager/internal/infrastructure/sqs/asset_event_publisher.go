package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type EventPublisher struct {
	producer    *sqs.Producer
	jobProducer *sqs.Producer
	logger      *logger.Logger
}

func NewEventPublisher(producer *sqs.Producer) *EventPublisher {
	return &EventPublisher{
		producer: producer,
		logger:   logger.WithService("sqs-event-publisher"),
	}
}

func NewEventPublisherWithJobProducer(producer *sqs.Producer, jobProducer *sqs.Producer) *EventPublisher {
	return &EventPublisher{
		producer:    producer,
		jobProducer: jobProducer,
		logger:      logger.WithService("sqs-event-publisher"),
	}
}

func (p *EventPublisher) PublishAssetCreated(ctx context.Context, a *domainasset.Asset) error {
	event := map[string]interface{}{
		"event":     "asset.created",
		"assetId":   a.ID().Value(),
		"slug":      a.Slug().Value(),
		"timestamp": a.CreatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishAssetUpdated(ctx context.Context, a *domainasset.Asset) error {
	event := map[string]interface{}{
		"event":     "asset.updated",
		"assetId":   a.ID().Value(),
		"slug":      a.Slug().Value(),
		"timestamp": a.UpdatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishAssetDeleted(ctx context.Context, a *domainasset.Asset) error {
	event := map[string]interface{}{
		"event":     "asset.deleted",
		"assetId":   a.ID().Value(),
		"slug":      a.Slug().Value(),
		"timestamp": a.UpdatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishAssetPublished(ctx context.Context, a *domainasset.Asset) error {
	event := map[string]interface{}{
		"event":     "asset.published",
		"assetId":   a.ID().Value(),
		"slug":      a.Slug().Value(),
		"timestamp": a.UpdatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishVideoAdded(ctx context.Context, a *domainasset.Asset, video *domainasset.Video) error {
	event := map[string]interface{}{
		"event":     "video.added",
		"assetId":   a.ID().Value(),
		"videoId":   video.ID(),
		"label":     video.Label(),
		"format":    video.Format(),
		"timestamp": video.CreatedAt(),
	}

	if err := p.publishEvent(ctx, event); err != nil {
		return err
	}

	if video.Format() == domainasset.VideoFormat(constants.VideoStreamingFormatRaw) {
		if err := p.triggerAnalyzeJob(ctx, a.ID().Value(), video.ID(), video.StorageLocation()); err != nil {
			p.logger.WithError(err).Error("Failed to trigger analyze job", "asset_id", a.ID().Value(), "video_id", video.ID())
		}
	} else if video.Format() == domainasset.VideoFormat(constants.VideoStreamingFormatHLS) || video.Format() == domainasset.VideoFormat(constants.VideoStreamingFormatDASH) {
		if err := p.triggerTranscodeJob(ctx, a, video.ID(), video.StorageLocation(), string(video.Format())); err != nil {
			p.logger.WithError(err).Error("Failed to trigger transcode job", "asset_id", a.ID().Value(), "video_id", video.ID(), "format", video.Format())
		}
	}

	return nil
}

func (p *EventPublisher) PublishVideoRemoved(ctx context.Context, a *domainasset.Asset, videoID string) error {
	event := map[string]interface{}{
		"event":     "video.removed",
		"assetId":   a.ID().Value(),
		"videoId":   videoID,
		"timestamp": a.UpdatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishVideoStatusUpdated(ctx context.Context, a *domainasset.Asset, videoID string, status domainasset.VideoStatus) error {
	event := map[string]interface{}{
		"event":     "video.status.updated",
		"assetId":   a.ID().Value(),
		"videoId":   videoID,
		"status":    status,
		"timestamp": a.UpdatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishImageAdded(ctx context.Context, a *domainasset.Asset, image domainasset.Image) error {
	event := map[string]interface{}{
		"event":     "image.added",
		"assetId":   a.ID().Value(),
		"imageId":   image.ID,
		"fileName":  image.FileName,
		"type":      image.Type,
		"timestamp": image.CreatedAt,
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishImageRemoved(ctx context.Context, a *domainasset.Asset, imageID string) error {
	event := map[string]interface{}{
		"event":     "image.removed",
		"assetId":   a.ID().Value(),
		"imageId":   imageID,
		"timestamp": a.UpdatedAt().Value(),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishBucketCreated(ctx context.Context, bucket *domainbucket.Bucket) error {
	event := map[string]interface{}{
		"event":     "bucket.created",
		"bucketId":  bucket.ID(),
		"name":      bucket.Name(),
		"key":       bucket.Key(),
		"timestamp": bucket.CreatedAt().Format(time.RFC3339),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishBucketUpdated(ctx context.Context, bucket *domainbucket.Bucket) error {
	event := map[string]interface{}{
		"event":     "bucket.updated",
		"bucketId":  bucket.ID(),
		"name":      bucket.Name(),
		"key":       bucket.Key(),
		"timestamp": bucket.UpdatedAt().Format(time.RFC3339),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishBucketDeleted(ctx context.Context, bucketID string) error {
	event := map[string]interface{}{
		"event":     "bucket.deleted",
		"bucketId":  bucketID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishAssetAddedToBucket(ctx context.Context, bucketID string, assetID string) error {
	event := map[string]interface{}{
		"event":     "bucket.asset.added",
		"bucketId":  bucketID,
		"assetId":   assetID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) PublishAssetRemovedFromBucket(ctx context.Context, bucketID string, assetID string) error {
	event := map[string]interface{}{
		"event":     "bucket.asset.removed",
		"bucketId":  bucketID,
		"assetId":   assetID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return p.publishEvent(ctx, event)
}

func (p *EventPublisher) publishEvent(ctx context.Context, event map[string]interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		p.logger.WithError(err).Error("Failed to marshal event", "event", event)
		return err
	}

	if err := p.producer.SendMessage(ctx, "domain_event", string(payload)); err != nil {
		p.logger.WithError(err).Error("Failed to publish event", "event", event)
		return err
	}

	p.logger.Debug("Event published successfully", "event", event)
	return nil
}

func (p *EventPublisher) triggerAnalyzeJob(ctx context.Context, assetID, videoID string, storageLocation domainasset.S3Object) error {
	if p.jobProducer == nil {
		p.logger.Warn("Job producer not configured, skipping analyze job", "asset_id", assetID, "video_id", videoID)
		return nil
	}

	input := fmt.Sprintf("s3://%s/%s", storageLocation.Bucket(), storageLocation.Key())

	payload := messages.JobPayload{
		JobType: "analyze",
		Input:   input,
		AssetID: assetID,
		VideoID: videoID,
	}

	err := p.jobProducer.SendMessage(ctx, messages.MessageTypeJob, payload)
	if err != nil {
		return fmt.Errorf("failed to send analyze job: %w", err)
	}

	p.logger.Info("Analyze job triggered successfully", "asset_id", assetID, "video_id", videoID, "input", input)
	return nil
}

func (p *EventPublisher) triggerTranscodeJob(ctx context.Context, asset *domainasset.Asset, videoID string, storageLocation domainasset.S3Object, format string) error {
	if p.jobProducer == nil {
		p.logger.Warn("Job producer not configured, skipping transcode job", "asset_id", asset.ID().Value(), "video_id", videoID)
		return nil
	}

	var input string
	if format == "hls" || format == "dash" {
		rawVideo := p.findRawVideo(asset)
		if rawVideo == nil {
			p.logger.Error("No raw video found for transcode job", "asset_id", asset.ID().Value(), "video_id", videoID, "format", format)
			return fmt.Errorf("no raw video found for transcode job")
		}
		input = fmt.Sprintf("s3://%s/%s", rawVideo.StorageLocation().Bucket(), rawVideo.StorageLocation().Key())
	} else {
		input = fmt.Sprintf("s3://%s/%s", storageLocation.Bucket(), storageLocation.Key())
	}

	outputBucket := "content-east"
	var outputKey string
	if format == "dash" {
		outputKey = fmt.Sprintf("%s/%s/%s/playlist.mpd", asset.ID().Value(), videoID, format)
	} else {
		outputKey = fmt.Sprintf("%s/%s/%s/playlist.%s", asset.ID().Value(), videoID, format, p.getOutputExtension(format))
	}

	payload := messages.JobPayload{
		JobType:      "transcode",
		Input:        input,
		AssetID:      asset.ID().Value(),
		VideoID:      videoID,
		Format:       format,
		OutputBucket: outputBucket,
		OutputKey:    outputKey,
	}

	err := p.jobProducer.SendMessage(ctx, messages.MessageTypeJob, payload)
	if err != nil {
		return fmt.Errorf("failed to send transcode job: %w", err)
	}

	p.logger.Info("Transcode job triggered successfully", "asset_id", asset.ID().Value(), "video_id", videoID, "format", format, "input", input, "output", fmt.Sprintf("s3://%s/%s", outputBucket, outputKey))
	return nil
}

func (p *EventPublisher) findRawVideo(asset *domainasset.Asset) *domainasset.Video {
	for _, video := range asset.Videos() {
		if video.Format() == domainasset.VideoFormat(constants.VideoStreamingFormatRaw) {
			return video
		}
	}
	return nil
}

func (p *EventPublisher) getOutputExtension(format string) string {
	switch format {
	case "hls":
		return "m3u8"
	case "dash":
		return "mpd"
	default:
		return "mp4"
	}
}
