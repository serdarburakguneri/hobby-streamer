package asset

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type AssetService interface {
	ListAssets(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	SearchAssets(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	GetAssetByID(ctx context.Context, id string) (*Asset, error)
	GetAssetBySlug(ctx context.Context, slug string) (*Asset, error)
	GetAssetsByIDs(ctx context.Context, ids []string) ([]Asset, error)
	CreateAsset(ctx context.Context, a *Asset) (*Asset, error)
	UpdateAsset(ctx context.Context, id string, a *Asset) error
	PatchAsset(ctx context.Context, id string, patch map[string]interface{}) error
	DeleteAsset(ctx context.Context, id string) error
	AddImage(ctx context.Context, id string, img *Image) error
	DeleteImage(ctx context.Context, id string, filename string) error
	AddVideo(ctx context.Context, id string, video *Video) error
	DeleteVideo(ctx context.Context, assetID string, videoID string) error
	HandleAnalyzeCompletion(ctx context.Context, payload map[string]interface{}) error
	HandleTranscodeCompletion(ctx context.Context, payload map[string]interface{}) error
	GetChildren(ctx context.Context, parentID string) ([]Asset, error)
	GetParent(ctx context.Context, childID string) (*Asset, error)
	GetAssetsByTypeAndGenre(ctx context.Context, assetType, genre string) ([]Asset, error)
}

type Service struct {
	Repo   AssetRepository
	Config *config.DynamicConfig
}

var _ AssetService = (*Service)(nil)

type S3File struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type DeleteFilesRequest struct {
	AssetID string   `json:"assetId"`
	Files   []S3File `json:"files"`
}

func generateID() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000_000))
	if err != nil {
		return "000000000"
	}
	return strconv.Itoa(int(n.Int64()))
}

func generateSlug(title string) string {
	if title == "" {
		return generateID()
	}

	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`[^a-z0-9-_]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`_+`).ReplaceAllString(slug, "_")
	slug = strings.Trim(slug, "-_")

	if slug == "" {
		return generateID()
	}

	return slug
}

func (s *Service) GetAssetByID(ctx context.Context, id string) (*Asset, error) {
	return s.Repo.GetAssetByID(ctx, id)
}

func (s *Service) GetAssetBySlug(ctx context.Context, slug string) (*Asset, error) {
	return s.Repo.GetAssetBySlug(ctx, slug)
}

func (s *Service) GetAssetsByIDs(ctx context.Context, ids []string) ([]Asset, error) {
	return s.Repo.GetAssetsByIDs(ctx, ids)
}

func (s *Service) ListAssets(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	return s.Repo.ListAssets(ctx, limit, lastKey)
}

func (s *Service) SearchAssets(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	return s.Repo.SearchAssets(ctx, query, limit, lastKey)
}

func (s *Service) CreateAsset(ctx context.Context, a *Asset) (*Asset, error) {
	log := logger.Get().WithService("asset-service")

	if a.ID != "" {
		log.Error("Asset ID should not be set during creation", "provided_id", a.ID)
		return nil, apperrors.NewValidationError("asset ID should not be set during creation", nil)
	}

	if err := s.validateAsset(a); err != nil {
		log.WithError(err).Error("Asset validation failed")
		return nil, err
	}

	a.ID = generateID()

	if a.Slug == "" {
		title := ""
		if a.Title != nil {
			title = *a.Title
		}
		a.Slug = generateSlug(title)
	}

	existingAsset, err := s.Repo.GetAssetBySlug(ctx, a.Slug)
	if err == nil && existingAsset != nil {
		log.Error("Slug already exists", "slug", a.Slug, "existing_asset_id", existingAsset.ID)
		return nil, apperrors.NewConflictError("slug already exists", nil)
	}

	if err := s.Repo.SaveAsset(ctx, a); err != nil {
		log.WithError(err).Error("Failed to save asset to repository", "asset_id", a.ID)
		return nil, apperrors.NewInternalError("failed to save asset to repository", err)
	}

	return a, nil
}

func (s *Service) UpdateAsset(ctx context.Context, id string, a *Asset) error {
	if a.ID != id {
		return apperrors.NewValidationError("asset ID mismatch", nil)
	}

	if err := s.validateAsset(a); err != nil {
		return err
	}

	return s.Repo.SaveAsset(ctx, a)
}

func (s *Service) PatchAsset(ctx context.Context, id string, patch map[string]interface{}) error {
	if _, exists := patch["slug"]; exists {
		return apperrors.NewValidationError("slug cannot be updated after creation", nil)
	}
	return s.Repo.PatchAsset(ctx, id, patch)
}

func (s *Service) DeleteAsset(ctx context.Context, id string) error {
	log := logger.Get().WithService("asset-service")

	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	var filesToDelete []S3File

	for _, video := range asset.Videos {
		filesToDelete = append(filesToDelete, S3File{
			Bucket: video.StorageLocation.Bucket,
			Key:    video.StorageLocation.Key,
		})
		if video.Thumbnail != nil && video.Thumbnail.StorageLocation != nil {
			filesToDelete = append(filesToDelete, S3File{
				Bucket: video.Thumbnail.StorageLocation.Bucket,
				Key:    video.Thumbnail.StorageLocation.Key,
			})
		}
	}

	if err := s.Repo.DeleteAsset(ctx, id); err != nil {
		return err
	}

	if len(filesToDelete) > 0 {
		go func() {
			bgCtx := context.Background()
			if err := s.deleteAssetFiles(bgCtx, asset); err != nil {
				log.WithError(err).Error("Failed to delete asset files", "asset_id", id)
			} else {
				log.Debug("Asset files deleted successfully", "asset_id", id, "file_count", len(filesToDelete))
			}
		}()
	}

	return nil
}

func (s *Service) deleteAssetFiles(ctx context.Context, asset *Asset) error {
	if len(asset.Videos) == 0 {
		return nil
	}

	deleteRequest := map[string]interface{}{
		"assetId": asset.ID,
		"files":   []string{},
	}

	for _, video := range asset.Videos {
		if video.StorageLocation.Bucket != "" && video.StorageLocation.Key != "" {
			filePath := fmt.Sprintf("%s/%s", video.StorageLocation.Bucket, video.StorageLocation.Key)
			deleteRequest["files"] = append(deleteRequest["files"].([]string), filePath)
		}
	}

	if len(deleteRequest["files"].([]string)) == 0 {
		return nil
	}

	lambdaURL := s.Config.GetStringFromComponent("lambda", "delete_files_endpoint")
	if lambdaURL == "" {
		return apperrors.NewInternalError("delete_files_endpoint not configured", nil)
	}

	requestBody, err := json.Marshal(deleteRequest)
	if err != nil {
		return apperrors.NewInternalError("failed to marshal delete request", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", lambdaURL, strings.NewReader(string(requestBody)))
	if err != nil {
		return apperrors.NewInternalError("failed to create delete request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return apperrors.NewExternalError("failed to call delete-files function", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperrors.NewExternalError(fmt.Sprintf("delete-files function returned status: %d", resp.StatusCode), nil)
	}

	return nil
}

func (s *Service) AddImage(ctx context.Context, id string, img *Image) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	for _, existing := range asset.Images {
		if existing.FileName == img.FileName {
			return apperrors.NewConflictError("image already exists for this asset", nil)
		}
	}

	asset.Images = append(asset.Images, *img)
	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) DeleteImage(ctx context.Context, id string, filename string) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	filtered := make([]Image, 0, len(asset.Images))
	for _, img := range asset.Images {
		if img.FileName != filename {
			filtered = append(filtered, img)
		}
	}

	asset.Images = filtered
	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) AddVideo(ctx context.Context, id string, video *Video) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
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
	err = s.Repo.SaveAsset(ctx, asset)
	if err != nil {
		return err
	}

	if video.Format == VideoFormatRaw {
		s.sendAnalyzeJob(ctx, asset.ID, video.ID, video.StorageLocation)
	}
	return nil
}

func (s *Service) sendAnalyzeJob(ctx context.Context, assetID string, videoID string, storageLocation S3Object) {
	log := logger.Get().WithService("asset-service")

	input := fmt.Sprintf("s3://%s/%s", storageLocation.Bucket, storageLocation.Key)
	payload := messages.AnalyzePayload{
		Input:   input,
		AssetID: assetID,
		VideoID: videoID,
	}

	analyzeJobsQueueURL := s.Config.GetStringFromComponent("sqs", "analyze_jobs_queue_url")
	if analyzeJobsQueueURL == "" {
		log.Error("Analyze jobs queue URL not configured", "asset_id", assetID, "video_id", videoID)
		return
	}

	producer, err := sqs.NewProducer(ctx, analyzeJobsQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create analyze SQS producer", "asset_id", assetID, "video_id", videoID)
		return
	}

	err = producer.SendMessage(ctx, messages.MessageTypeAnalyze, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze job", "asset_id", assetID, "video_id", videoID, "input", input)
	} else {
		log.Info("Analyze job sent successfully", "asset_id", assetID, "video_id", videoID, "input", input)
	}
}

func (s *Service) DeleteVideo(ctx context.Context, assetID string, videoID string) error {
	asset, err := s.Repo.GetAssetByID(ctx, assetID)
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
	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) GetChildren(ctx context.Context, parentID string) ([]Asset, error) {
	return s.Repo.GetChildren(ctx, parentID)
}

func (s *Service) GetParent(ctx context.Context, childID string) (*Asset, error) {
	return s.Repo.GetParent(ctx, childID)
}

func (s *Service) GetAssetsByTypeAndGenre(ctx context.Context, assetType, genre string) ([]Asset, error) {
	return s.Repo.GetAssetsByTypeAndGenre(ctx, assetType, genre)
}

func (s *Service) HandleAnalyzeCompletion(ctx context.Context, payload map[string]interface{}) error {
	log := logger.Get().WithService("asset-service")

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

	asset, err := s.Repo.GetAssetByID(ctx, analyzePayload.AssetID)
	if err != nil {
		log.WithError(err).Error("Failed to get asset for analyze completion", "asset_id", analyzePayload.AssetID)
		return err
	}

	for i, video := range asset.Videos {
		if video.ID == analyzePayload.VideoID {
			asset.Videos[i].Status = status
			asset.Videos[i].UpdatedAt = time.Now()
			break
		}
	}

	err = s.Repo.SaveAsset(ctx, asset)
	if err != nil {
		log.WithError(err).Error("Failed to save asset after analyze completion", "asset_id", analyzePayload.AssetID, "video_id", analyzePayload.VideoID, "status", status)
		return err
	}

	log.Info("Analyze completion processed successfully", "asset_id", analyzePayload.AssetID, "video_id", analyzePayload.VideoID, "success", analyzePayload.Success, "status", status)
	return nil
}

func (s *Service) HandleTranscodeCompletion(ctx context.Context, payload map[string]interface{}) error {
	log := logger.Get().WithService("asset-service")

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

	if transcodePayload.Success {
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
		log.Debug("CDN prefix lookup", "bucket", transcodePayload.Bucket, "cdn_prefix", cdnPrefix)
		if cdnPrefix != "" {
			playURL := cdnPrefix + "/" + transcodePayload.Key
			video.StreamInfo = &StreamInfo{
				CdnPrefix: &cdnPrefix,
				PlayURL:   &playURL,
			}
			log.Debug("Set CDN prefix and play URL for video", "asset_id", transcodePayload.AssetID, "bucket", transcodePayload.Bucket, "cdn_prefix", cdnPrefix, "play_url", playURL)
		} else {
			log.Warn("No CDN prefix found for bucket", "bucket", transcodePayload.Bucket)
		}

		err = s.AddVideo(ctx, transcodePayload.AssetID, video)
		if err != nil {
			log.WithError(err).Error("Failed to add transcoded video", "asset_id", transcodePayload.AssetID, "video_id", video.ID)
			return err
		}

		log.Info("Created new transcoded video", "asset_id", transcodePayload.AssetID, "video_id", video.ID, "format", transcodePayload.Format, "s3_location", s3Object)
	} else {
		log.Error("Transcode failed", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID, "error", transcodePayload.Error)
	}

	log.Info("Transcode completion processed successfully", "asset_id", transcodePayload.AssetID, "video_id", transcodePayload.VideoID, "format", transcodePayload.Format, "success", transcodePayload.Success, "status", status)
	return nil
}

func (s *Service) getCDNPrefixForBucket(bucket string) string {
	switch bucket {
	case "hls-storage":
		return s.Config.GetStringFromComponent("cdn", "hls_prefix")
	case "dash-storage":
		return s.Config.GetStringFromComponent("cdn", "dash_prefix")
	case "thumbnails-storage":
		return s.Config.GetStringFromComponent("cdn", "thumbnails_prefix")
	default:
		return ""
	}
}

func (s *Service) validateAsset(a *Asset) error {
	log := logger.Get().WithService("asset-service")

	if a.Slug != "" {
		if !isValidSlug(a.Slug) {
			log.Error("Invalid slug format", "slug", a.Slug)
			return apperrors.NewValidationError("invalid slug format", nil)
		}
	}

	if a.Type != nil {
		validTypes := []string{
			AssetTypeMovie, AssetTypeSeries, AssetTypeSeason, AssetTypeEpisode,
			AssetTypeDocumentary, AssetTypeMusic, AssetTypePodcast,
			AssetTypeTrailer, AssetTypeBehindTheScenes, AssetTypeInterview,
		}
		if !contains(validTypes, *a.Type) {
			log.Error("Invalid asset type", "type", *a.Type, "valid_types", validTypes)
			return apperrors.NewValidationError("invalid type value", nil)
		}
	}

	if a.Genre != nil {
		validGenres := []string{
			AssetGenreAction, AssetGenreDrama, AssetGenreComedy, AssetGenreHorror,
			AssetGenreSciFi, AssetGenreRomance, AssetGenreThriller, AssetGenreFantasy,
			AssetGenreDocumentary, AssetGenreMusic, AssetGenreNews,
			AssetGenreSports, AssetGenreKids, AssetGenreEducational,
		}
		if !contains(validGenres, *a.Genre) {
			log.Error("Invalid primary genre", "genre", *a.Genre, "valid_genres", validGenres)
			return apperrors.NewValidationError("invalid genre value", nil)
		}
	}

	if len(a.Genres) > 0 {
		validGenres := []string{
			AssetGenreAction, AssetGenreDrama, AssetGenreComedy, AssetGenreHorror,
			AssetGenreSciFi, AssetGenreRomance, AssetGenreThriller, AssetGenreFantasy,
			AssetGenreDocumentary, AssetGenreMusic, AssetGenreNews,
			AssetGenreSports, AssetGenreKids, AssetGenreEducational,
		}
		for _, genre := range a.Genres {
			if !contains(validGenres, genre) {
				log.Error("Invalid additional genre", "genre", genre, "valid_genres", validGenres)
				return apperrors.NewValidationError("invalid genre value in genres array", nil)
			}
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-z0-9-_]+$`, slug)
	return matched
}

func NewService(repo AssetRepository, cfg *config.DynamicConfig) *Service {
	return &Service{
		Repo:   repo,
		Config: cfg,
	}
}

func NewServiceWithSQS(repo AssetRepository, sqsProducer *sqs.Producer, cfg *config.DynamicConfig) *Service {
	return &Service{
		Repo:   repo,
		Config: cfg,
	}
}
