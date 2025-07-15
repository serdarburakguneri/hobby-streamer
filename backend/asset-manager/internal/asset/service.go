package asset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	UpdateVideo(ctx context.Context, assetID string, videoID string, video *Video) error
	DeleteVideo(ctx context.Context, assetID string, videoID string) error
	UpdateVideoStatus(ctx context.Context, assetID string, videoID string, status string) error
	UpdateVideoCDN(ctx context.Context, assetID string, videoID string, cdnPrefix string) error
	HandleAnalyzeCompletion(ctx context.Context, payload map[string]interface{}) error
	HandleTranscodeCompletion(ctx context.Context, payload map[string]interface{}) error
	GetChildren(ctx context.Context, parentID string) ([]Asset, error)
	GetParent(ctx context.Context, childID string) (*Asset, error)
	GetAssetsByTypeAndGenre(ctx context.Context, assetType, genre string) ([]Asset, error)
}

type Service struct {
	Repo        AssetRepository
	SQSProducer *sqs.Producer
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
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(1_000_000_000))
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
	if len(ids) == 0 {
		return []Asset{}, nil
	}

	assets := make([]Asset, 0, len(ids))
	for _, id := range ids {
		asset, err := s.Repo.GetAssetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get asset %s: %w", id, err)
		}
		assets = append(assets, *asset)
	}

	return assets, nil
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
		return nil, errors.New(ErrIDShouldNotBeSet)
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
		return nil, errors.New("slug already exists")
	}

	if err := s.Repo.SaveAsset(ctx, a); err != nil {
		log.WithError(err).Error("Failed to save asset to repository", "asset_id", a.ID)
		return nil, err
	}

	return a, nil
}

func (s *Service) UpdateAsset(ctx context.Context, id string, a *Asset) error {
	if a.ID != id {
		return errors.New(ErrIDMismatch)
	}

	if err := s.validateAsset(a); err != nil {
		return err
	}

	return s.Repo.SaveAsset(ctx, a)
}

func (s *Service) PatchAsset(ctx context.Context, id string, patch map[string]interface{}) error {
	if _, exists := patch["slug"]; exists {
		return errors.New("slug cannot be updated after creation")
	}
	return s.Repo.PatchAsset(ctx, id, patch)
}

func (s *Service) DeleteAsset(ctx context.Context, id string) error {
	log := logger.Get().WithService("asset-service")

	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
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
			if err := s.deleteAssetFiles(bgCtx, id, filesToDelete); err != nil {
				log.WithError(err).Error("Failed to delete asset files", "asset_id", id)
			} else {
				log.Debug("Asset files deleted successfully", "asset_id", id, "file_count", len(filesToDelete))
			}
		}()
	}

	return nil
}

func (s *Service) deleteAssetFiles(ctx context.Context, assetID string, files []S3File) error {

	request := DeleteFilesRequest{
		AssetID: assetID,
		Files:   files,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal delete request: %w", err)
	}

	lambdaEndpoint := os.Getenv("DELETE_FILES_LAMBDA_ENDPOINT")
	if lambdaEndpoint == "" {
		lambdaEndpoint = "http://localhost:4566/2015-03-31/functions/delete-files/invocations"
	}

	resp, err := http.Post(
		lambdaEndpoint,
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to call delete-files function: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete-files function returned status: %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) AddImage(ctx context.Context, id string, img *Image) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	for _, existing := range asset.Images {
		if existing.FileName == img.FileName {
			return errors.New(ErrImageExists)
		}
	}

	asset.Images = append(asset.Images, *img)
	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) DeleteImage(ctx context.Context, id string, filename string) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
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
		return err
	}

	if asset.Videos == nil {
		asset.Videos = make([]Video, 0)
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

	if s.SQSProducer != nil && video.Format == VideoFormatRaw {
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

	err := s.SQSProducer.SendMessage(ctx, messages.MessageTypeAnalyze, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze job", "asset_id", assetID, "video_id", videoID, "input", input)
	} else {
		log.Info("Analyze job sent successfully", "asset_id", assetID, "video_id", videoID, "input", input)
	}
}

func (s *Service) DeleteVideo(ctx context.Context, assetID string, videoID string) error {
	asset, err := s.Repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
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

// --- Helper Functions ---

func getVideoByID(asset *Asset, videoID string) *Video {
	for i := range asset.Videos {
		if asset.Videos[i].ID == videoID {
			return &asset.Videos[i]
		}
	}
	return nil
}

func (s *Service) setCDNPrefixIfAvailable(video *Video, bucket string) {
	cdnPrefix := s.getCDNPrefixForBucket(bucket)
	if cdnPrefix != "" {
		if video.StreamInfo == nil {
			video.StreamInfo = &StreamInfo{}
		}
		video.StreamInfo.CdnPrefix = &cdnPrefix
		playURL := cdnPrefix + "/" + video.StorageLocation.Key
		video.StreamInfo.PlayURL = &playURL
	}
}

// --- AssetService Methods ---

func (s *Service) UpdateVideo(ctx context.Context, assetID string, videoID string, video *Video) error {
	asset, err := s.Repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}

	existingVideo := getVideoByID(asset, videoID)
	if existingVideo == nil {
		return errors.New("video not found")
	}

	video.UpdatedAt = time.Now()
	*existingVideo = *video
	existingVideo.ID = videoID

	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) UpdateVideoStatus(ctx context.Context, assetID string, videoID string, status string) error {
	asset, err := s.Repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}

	video := getVideoByID(asset, videoID)
	if video == nil {
		return errors.New("video not found")
	}

	video.Status = status
	video.UpdatedAt = time.Now()

	if video.Format == VideoFormatHLS || video.Format == VideoFormatDASH {
		s.setCDNPrefixIfAvailable(video, video.StorageLocation.Bucket)
	}

	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) UpdateVideoCDN(ctx context.Context, assetID string, videoID string, cdnPrefix string) error {
	asset, err := s.Repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}

	video := getVideoByID(asset, videoID)
	if video == nil {
		return errors.New("video not found")
	}

	if video.StreamInfo == nil {
		video.StreamInfo = &StreamInfo{}
	}
	video.StreamInfo.CdnPrefix = &cdnPrefix
	video.UpdatedAt = time.Now()

	return s.Repo.SaveAsset(ctx, asset)
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

	err = s.UpdateVideoStatus(ctx, analyzePayload.AssetID, analyzePayload.VideoID, status)
	if err != nil {
		log.WithError(err).Error("Failed to update video status after analyze completion", "asset_id", analyzePayload.AssetID, "video_id", analyzePayload.VideoID, "status", status)
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
		return os.Getenv("HLS_CDN_PREFIX")
	case "dash-storage":
		return os.Getenv("DASH_CDN_PREFIX")
	case "thumbnails-storage":
		return os.Getenv("THUMBNAILS_CDN_PREFIX")
	default:
		return ""
	}
}

func (s *Service) validateAsset(a *Asset) error {
	log := logger.Get().WithService("asset-service")

	if a.Slug != "" {
		if !isValidSlug(a.Slug) {
			log.Error("Invalid slug format", "slug", a.Slug)
			return errors.New("invalid slug format")
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
			return errors.New("invalid type value")
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
			return errors.New("invalid genre value")
		}
	}

	if a.Genres != nil && len(a.Genres) > 0 {
		validGenres := []string{
			AssetGenreAction, AssetGenreDrama, AssetGenreComedy, AssetGenreHorror,
			AssetGenreSciFi, AssetGenreRomance, AssetGenreThriller, AssetGenreFantasy,
			AssetGenreDocumentary, AssetGenreMusic, AssetGenreNews,
			AssetGenreSports, AssetGenreKids, AssetGenreEducational,
		}
		for _, genre := range a.Genres {
			if !contains(validGenres, genre) {
				log.Error("Invalid additional genre", "genre", genre, "valid_genres", validGenres)
				return errors.New("invalid genre value in genres array")
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

func NewService(repo AssetRepository) *Service {
	return &Service{Repo: repo}
}

func NewServiceWithSQS(repo AssetRepository, sqsProducer *sqs.Producer) *Service {
	return &Service{
		Repo:        repo,
		SQSProducer: sqsProducer,
	}
}
