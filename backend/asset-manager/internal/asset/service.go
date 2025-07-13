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
	AddVideo(ctx context.Context, id string, videoType VideoType, video *Video) error
	DeleteVideo(ctx context.Context, id string, videoType VideoType) error
	UpdateVideoStatus(ctx context.Context, id string, videoType VideoType, status string) error
	AddVideoVariant(ctx context.Context, id string, videoType VideoType, variant string, videoVariant *VideoVariant) error
	UpdateVideoVariant(ctx context.Context, id string, videoType VideoType, variant string, videoVariant *VideoVariant) error
	DeleteVideoVariant(ctx context.Context, id string, videoType VideoType, variant string) error
	UpdateVideoVariantStatus(ctx context.Context, id string, videoType VideoType, variant string, status string) error
	HandleStatusUpdateMessage(ctx context.Context, messageType string, payload map[string]interface{}) error
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
		if video.Raw != nil {
			filesToDelete = append(filesToDelete, S3File{
				Bucket: video.Raw.StorageLocation.Bucket,
				Key:    video.Raw.StorageLocation.Key,
			})
		}
		if video.HLS != nil {
			filesToDelete = append(filesToDelete, S3File{
				Bucket: video.HLS.StorageLocation.Bucket,
				Key:    video.HLS.StorageLocation.Key,
			})
		}
		if video.DASH != nil {
			filesToDelete = append(filesToDelete, S3File{
				Bucket: video.DASH.StorageLocation.Bucket,
				Key:    video.DASH.StorageLocation.Key,
			})
		}
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

func (s *Service) AddVideo(ctx context.Context, id string, videoType VideoType, video *Video) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	if asset.Videos == nil {
		asset.Videos = make([]Video, 0)
	}

	for i, existingVideo := range asset.Videos {
		if existingVideo.Type == videoType {
			video.Type = videoType
			if video.Raw != nil && video.Raw.Status == "" {
				video.Raw.Status = VideoStatusPending
			}
			asset.Videos[i] = *video
			err = s.Repo.SaveAsset(ctx, asset)
			if err != nil {
				return err
			}

			if s.SQSProducer != nil && video.Raw != nil {
				s.sendAnalyzeJob(ctx, asset.ID, videoType, video.Raw.StorageLocation)
			}
			return nil
		}
	}

	video.Type = videoType
	if video.Raw != nil && video.Raw.Status == "" {
		video.Raw.Status = VideoStatusPending
	}
	asset.Videos = append(asset.Videos, *video)
	err = s.Repo.SaveAsset(ctx, asset)
	if err != nil {
		return err
	}

	if s.SQSProducer != nil && video.Raw != nil {
		s.sendAnalyzeJob(ctx, asset.ID, videoType, video.Raw.StorageLocation)
	}
	return nil
}

func (s *Service) sendAnalyzeJob(ctx context.Context, assetID string, videoType VideoType, storageLocation S3Object) {
	log := logger.Get().WithService("asset-service")

	input := fmt.Sprintf("s3://%s/%s", storageLocation.Bucket, storageLocation.Key)
	payload := map[string]interface{}{
		"input":     input,
		"assetId":   assetID,
		"videoType": string(videoType),
	}

	err := s.SQSProducer.SendMessage(ctx, "analyze", payload)
	if err != nil {
		log.WithError(err).Error("Failed to send analyze job", "asset_id", assetID, "input", input)
	} else {
		log.Info("Analyze job sent successfully", "asset_id", assetID, "input", input)
	}
}

func (s *Service) DeleteVideo(ctx context.Context, id string, videoType VideoType) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	filtered := make([]Video, 0, len(asset.Videos))
	for _, video := range asset.Videos {
		if video.Type != videoType {
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

func (s *Service) UpdateVideoStatus(ctx context.Context, id string, videoType VideoType, status string) error {
	return s.UpdateVideoVariantStatus(ctx, id, videoType, VideoVariantRaw, status)
}

func (s *Service) AddVideoVariant(ctx context.Context, id string, videoType VideoType, variant string, videoVariant *VideoVariant) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	for i, video := range asset.Videos {
		if video.Type == videoType {
			switch variant {
			case VideoVariantRaw:
				asset.Videos[i].Raw = videoVariant
			case VideoVariantHLS:
				asset.Videos[i].HLS = videoVariant
			case VideoVariantDASH:
				asset.Videos[i].DASH = videoVariant
			default:
				return errors.New("invalid variant")
			}
			return s.Repo.SaveAsset(ctx, asset)
		}
	}

	return errors.New("video not found")
}

func (s *Service) UpdateVideoVariant(ctx context.Context, id string, videoType VideoType, variant string, videoVariant *VideoVariant) error {
	return s.AddVideoVariant(ctx, id, videoType, variant, videoVariant)
}

func (s *Service) DeleteVideoVariant(ctx context.Context, id string, videoType VideoType, variant string) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	for i, video := range asset.Videos {
		if video.Type == videoType {
			switch variant {
			case VideoVariantRaw:
				asset.Videos[i].Raw = nil
			case VideoVariantHLS:
				asset.Videos[i].HLS = nil
			case VideoVariantDASH:
				asset.Videos[i].DASH = nil
			default:
				return errors.New("invalid variant")
			}
			return s.Repo.SaveAsset(ctx, asset)
		}
	}

	return errors.New("video not found")
}

func (s *Service) UpdateVideoVariantStatus(ctx context.Context, id string, videoType VideoType, variant string, status string) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	for i, video := range asset.Videos {
		if video.Type == videoType {
			switch variant {
			case VideoVariantRaw:
				if asset.Videos[i].Raw != nil {
					asset.Videos[i].Raw.Status = status
				}
			case VideoVariantHLS:
				if asset.Videos[i].HLS != nil {
					asset.Videos[i].HLS.Status = status
				}
			case VideoVariantDASH:
				if asset.Videos[i].DASH != nil {
					asset.Videos[i].DASH.Status = status
				}
			default:
				return errors.New("invalid variant")
			}
			return s.Repo.SaveAsset(ctx, asset)
		}
	}

	return errors.New("video not found")
}

func (s *Service) HandleStatusUpdateMessage(ctx context.Context, messageType string, payload map[string]interface{}) error {
	log := logger.Get().WithService("asset-service")

	if messageType != "update-video-variant-status" {
		return nil
	}

	assetID, ok := payload["assetId"].(string)
	if !ok {
		log.Error("Invalid assetId in status update message")
		return errors.New("invalid assetId")
	}

	videoTypeStr, ok := payload["videoType"].(string)
	if !ok {
		log.Error("Invalid videoType in status update message")
		return errors.New("invalid videoType")
	}

	variant, ok := payload["variant"].(string)
	if !ok {
		log.Error("Invalid variant in status update message")
		return errors.New("invalid variant")
	}

	status, ok := payload["status"].(string)
	if !ok {
		log.Error("Invalid status in status update message")
		return errors.New("invalid status")
	}

	var videoType VideoType
	switch videoTypeStr {
	case "main":
		videoType = VideoTypeMain
	case "trailer":
		videoType = VideoTypeTrailer
	case "behind_the_scenes":
		videoType = VideoTypeBehind
	case "interview":
		videoType = VideoTypeInterview
	default:
		log.Error("Unknown video type in status update message", "video_type", videoTypeStr)
		return errors.New("unknown video type")
	}

	log.Info("Processing status update", "asset_id", assetID, "video_type", videoType, "variant", variant, "status", status)

	err := s.UpdateVideoVariantStatus(ctx, assetID, videoType, variant, status)
	if err != nil {
		log.WithError(err).Error("Failed to update video variant status", "asset_id", assetID, "video_type", videoType, "variant", variant, "status", status)
		return err
	}

	log.Info("Status update processed successfully", "asset_id", assetID, "video_type", videoType, "variant", variant, "status", status)
	return nil
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
