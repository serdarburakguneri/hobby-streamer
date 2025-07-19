package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type VideoService struct {
	repo   AssetRepository
	config *config.DynamicConfig
}

func NewVideoService(repo AssetRepository, config *config.DynamicConfig) *VideoService {
	return &VideoService{
		repo:   repo,
		config: config,
	}
}

func (s *VideoService) AddVideo(ctx context.Context, assetID string, video *Video) error {
	asset, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	for i, existingVideo := range asset.Videos {
		if existingVideo.Format == video.Format {
			if existingVideo.Status == VideoStatusFailed {
				asset.Videos[i].Status = VideoStatusPending
				asset.Videos[i].UpdatedAt = time.Now()
				err = s.repo.SaveAsset(ctx, asset)
				if err != nil {
					return err
				}

				if video.Format == VideoFormatHLS || video.Format == VideoFormatDASH {
					jobSent := s.sendTranscodeJob(ctx, asset.ID, existingVideo.ID, existingVideo.StorageLocation, string(video.Format), existingVideo)
					if !jobSent {
						asset.Videos[i].Status = VideoStatusFailed
						asset.Videos[i].UpdatedAt = time.Now()
						s.repo.SaveAsset(ctx, asset)
					}
				}
				return nil
			}
			if existingVideo.Status == VideoStatusReady {
				return nil
			}

			return apperrors.NewConflictError("video with this format is already being processed", nil)
		}
	}

	if video.ID == "" {
		video.ID = generateID()
	}

	if video.Status == "" {
		video.Status = VideoStatusPending
	}

	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()

	asset.Videos = append(asset.Videos, *video)
	err = s.repo.SaveAsset(ctx, asset)
	if err != nil {
		return err
	}

	if video.Format == VideoFormatRaw {
		jobSent := s.sendAnalyzeJob(ctx, asset.ID, video.ID, video.StorageLocation)
		if !jobSent {
			for i, v := range asset.Videos {
				if v.ID == video.ID {
					asset.Videos[i].Status = VideoStatusFailed
					asset.Videos[i].UpdatedAt = time.Now()
					s.repo.SaveAsset(ctx, asset)
					break
				}
			}
		}
	} else if video.Format == VideoFormatHLS || video.Format == VideoFormatDASH {
		jobSent := s.sendTranscodeJob(ctx, asset.ID, video.ID, video.StorageLocation, string(video.Format), *video)
		if !jobSent {
			for i, v := range asset.Videos {
				if v.ID == video.ID {
					asset.Videos[i].Status = VideoStatusFailed
					asset.Videos[i].UpdatedAt = time.Now()
					s.repo.SaveAsset(ctx, asset)
					break
				}
			}
		}
	}
	return nil
}

func (s *VideoService) DeleteVideo(ctx context.Context, assetID string, videoID string) error {
	asset, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	filtered := make([]Video, 0, len(asset.Videos))
	for _, video := range asset.Videos {
		if video.ID != videoID {
			filtered = append(filtered, video)
		}
	}

	asset.Videos = filtered
	return s.repo.SaveAsset(ctx, asset)
}

func (s *VideoService) HandleAnalyzeCompletion(ctx context.Context, payload map[string]interface{}) error {
	log := logger.Get().WithService("video-service")

	var analyzePayload messages.AnalyzeCompletionPayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Failed to marshal analyze completion payload")
		return err
	}

	if err := json.Unmarshal(payloadBytes, &analyzePayload); err != nil {
		log.WithError(err).Error("Failed to unmarshal analyze completion payload")
		return err
	}

	var status string
	if analyzePayload.Success {
		status = VideoStatusReady
	} else {
		status = VideoStatusFailed
	}

	log.Info("Processing analyze completion", "asset_id", analyzePayload.AssetID, "video_id", analyzePayload.VideoID, "success", analyzePayload.Success, "status", status)

	asset, err := s.repo.GetAssetByID(ctx, analyzePayload.AssetID)
	if err != nil {
		log.WithError(err).Error("Failed to get asset for analyze completion", "asset_id", analyzePayload.AssetID)
		return err
	}

	for i, video := range asset.Videos {
		if video.ID == analyzePayload.VideoID {
			asset.Videos[i].Status = status
			asset.Videos[i].UpdatedAt = time.Now()

			if analyzePayload.Success {
				asset.Videos[i].Width = analyzePayload.Width
				asset.Videos[i].Height = analyzePayload.Height
				asset.Videos[i].Duration = analyzePayload.Duration
				asset.Videos[i].Bitrate = analyzePayload.Bitrate
				asset.Videos[i].Codec = analyzePayload.Codec
				asset.Videos[i].Size = analyzePayload.Size
				asset.Videos[i].ContentType = analyzePayload.ContentType
			}
			break
		}
	}

	err = s.repo.SaveAsset(ctx, asset)
	if err != nil {
		log.WithError(err).Error("Failed to save asset after analyze completion", "asset_id", analyzePayload.AssetID, "video_id", analyzePayload.VideoID, "status", status)
		return err
	}

	log.Info("Analyze completion processed successfully", "asset_id", analyzePayload.AssetID, "video_id", analyzePayload.VideoID, "success", analyzePayload.Success, "status", status)
	return nil
}

func (s *VideoService) HandleTranscodeCompletion(ctx context.Context, payload map[string]interface{}) error {
	log := logger.Get().WithService("video-service")

	var transcodePayload messages.TranscodeCompletionPayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Failed to marshal transcode completion payload")
		return err
	}

	if err := json.Unmarshal(payloadBytes, &transcodePayload); err != nil {
		log.WithError(err).Error("Failed to unmarshal transcode completion payload")
		return err
	}

	var status string
	if transcodePayload.Success {
		status = VideoStatusReady
	} else {
		status = VideoStatusFailed
	}

	log.Info("Processing transcode completion", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID, "format", transcodePayload.Format, "success", transcodePayload.Success, "status", status)

	asset, err := s.repo.GetAssetByID(ctx, transcodePayload.AssetID)
	if err != nil {
		log.WithError(err).Error("Failed to get asset for transcode completion", "asset_id", transcodePayload.AssetID)
		return err
	}

	videoUpdated := false
	for i, video := range asset.Videos {
		if video.Format == VideoFormat(transcodePayload.Format) {
			asset.Videos[i].Status = status
			asset.Videos[i].UpdatedAt = time.Now()

			if transcodePayload.Success {
				asset.Videos[i].StorageLocation = S3Object{
					Bucket: transcodePayload.Bucket,
					Key:    transcodePayload.Key,
					URL:    transcodePayload.URL,
				}

				if transcodePayload.Format == "dash" {
					asset.Videos[i].ContentType = "application/dash+xml"
				} else {
					asset.Videos[i].ContentType = "application/x-mpegURL"
				}

				cdnPrefix := s.getCDNPrefixForBucket(transcodePayload.Bucket)
				if cdnPrefix != "" {
					url := cdnPrefix + "/" + transcodePayload.Key
					asset.Videos[i].StreamInfo = &StreamInfo{
						CdnPrefix: &cdnPrefix,
						URL:       &url,
					}
				}
			}

			videoUpdated = true
			break
		}
	}

	if !videoUpdated && transcodePayload.Success {
		s3Object := S3Object{
			Bucket: transcodePayload.Bucket,
			Key:    transcodePayload.Key,
			URL:    transcodePayload.URL,
		}

		video := &Video{
			ID:              generateID(),
			Type:            VideoTypeMain,
			Format:          VideoFormat(transcodePayload.Format),
			StorageLocation: s3Object,
			Status:          status,
			ContentType:     "application/x-mpegURL",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if transcodePayload.Format == "dash" {
			video.ContentType = "application/dash+xml"
		}

		cdnPrefix := s.getCDNPrefixForBucket(transcodePayload.Bucket)
		if cdnPrefix != "" {
			url := cdnPrefix + "/" + transcodePayload.Key
			video.StreamInfo = &StreamInfo{
				CdnPrefix: &cdnPrefix,
				URL:       &url,
			}
		}

		asset.Videos = append(asset.Videos, *video)
		log.Info("Created new transcoded video", "asset_id", transcodePayload.AssetID, "video_id", video.ID, "format", transcodePayload.Format)
	} else if !videoUpdated && !transcodePayload.Success {
		log.Error("Transcode failed and no existing video found to update", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID, "error", transcodePayload.Error)
	} else if videoUpdated && !transcodePayload.Success {
		log.Error("Transcode failed, updated existing video status", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID, "error", transcodePayload.Error)
	}

	err = s.repo.SaveAsset(ctx, asset)
	if err != nil {
		log.WithError(err).Error("Failed to save asset after transcode completion", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID)
		return err
	}

	log.Info("Transcode completion processed successfully", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID, "format", transcodePayload.Format, "success", transcodePayload.Success, "status", status)
	return nil
}

func (s *VideoService) sendAnalyzeJob(ctx context.Context, assetID string, videoID string, storageLocation S3Object) bool {
	log := logger.Get().WithService("video-service")

	input := fmt.Sprintf("s3://%s/%s", storageLocation.Bucket, storageLocation.Key)
	payload := messages.AnalyzePayload{
		Input:   input,
		AssetID: assetID,
		VideoID: videoID,
	}

	analyzeJobsQueueURL := s.config.GetStringFromComponent("sqs", "analyze_jobs_queue_url")
	if analyzeJobsQueueURL == "" {
		log.Error("Analyze jobs queue URL not configured", "asset_id", assetID, "video_id", videoID)
		return false
	}

	producer, err := sqs.NewProducer(ctx, analyzeJobsQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create analyze SQS producer", "asset_id", assetID, "video_id", videoID)
		return false
	}

	err = producer.SendMessage(ctx, messages.MessageTypeAnalyze, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze job", "asset_id", assetID, "video_id", videoID, "input", input)
		return false
	}

	log.Info("Analyze job sent successfully", "asset_id", assetID, "video_id", videoID, "input", input)
	return true
}

func (s *VideoService) sendTranscodeJob(ctx context.Context, assetID string, videoID string, storageLocation S3Object, format string, video Video) bool {
	log := logger.Get().WithService("video-service")

	input := fmt.Sprintf("s3://%s/%s", storageLocation.Bucket, storageLocation.Key)

	var queueURL string
	var messageType string

	switch format {
	case "hls":
		queueURL = s.config.GetStringFromComponent("sqs", "hls_queue_url")
		messageType = messages.MessageTypeTranscodeHLS
	case "dash":
		queueURL = s.config.GetStringFromComponent("sqs", "dash_queue_url")
		messageType = messages.MessageTypeTranscodeDASH
	default:
		log.Error("Invalid format for transcode job", "format", format, "asset_id", assetID, "video_id", videoID)
		return false
	}

	if queueURL == "" {
		log.Error("Transcode queue URL not configured", "format", format, "asset_id", assetID, "video_id", videoID)
		return false
	}

	outputBucket := "content-east"
	quality := getQualityFromDimensions(video.Width, video.Height)

	var outputKey string
	var outputFileName string

	if format == "dash" {
		outputKey = fmt.Sprintf("%s/%s/%s/manifest.mpd", assetID, format, quality)
		outputFileName = "manifest.mpd"
	} else {
		outputKey = fmt.Sprintf("%s/%s/%s/playlist.%s", assetID, format, quality, getFileExtension(format))
		outputFileName = fmt.Sprintf("playlist.%s", getFileExtension(format))
	}

	payload := messages.TranscodePayload{
		Input:          input,
		AssetID:        assetID,
		VideoID:        videoID,
		Format:         format,
		Quality:        quality,
		OutputBucket:   outputBucket,
		OutputKey:      outputKey,
		OutputFileName: outputFileName,
	}

	producer, err := sqs.NewProducer(ctx, queueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create transcode SQS producer", "asset_id", assetID, "video_id", videoID, "format", format)
		return false
	}

	err = producer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode job", "asset_id", assetID, "video_id", videoID, "format", format, "input", input)
		return false
	}

	log.Info("Transcode job sent successfully", "asset_id", assetID, "video_id", videoID, "format", format, "input", input)
	return true
}

func (s *VideoService) getCDNPrefixForBucket(bucket string) string {
	switch bucket {
	case "content-east", "content-west":
		return s.config.GetStringFromComponent("cdn", "prefix")
	default:
		return ""
	}
}
