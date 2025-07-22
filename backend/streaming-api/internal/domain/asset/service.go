package asset

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type AssetPublishingService struct{}

func NewAssetPublishingService() *AssetPublishingService {
	return &AssetPublishingService{}
}

func (s *AssetPublishingService) GetPublishStatus(asset *Asset) constants.PublishStatus {
	if asset == nil {
		return constants.PublishStatusInvalid
	}

	if !asset.IsReady() {
		return constants.PublishStatusNotReady
	}

	if asset.PublishRule() == nil {
		return constants.PublishStatusNotConfigured
	}

	now := time.Now().UTC()

	if asset.PublishRule().PublishAt() != nil && now.Before(*asset.PublishRule().PublishAt()) {
		return constants.PublishStatusScheduled
	}

	if asset.PublishRule().UnpublishAt() != nil && now.After(*asset.PublishRule().UnpublishAt()) {
		return constants.PublishStatusExpired
	}

	if asset.IsPublished() {
		return constants.PublishStatusPublished
	}

	return constants.PublishStatusDraft
}

type AssetStreamingService struct{}

func NewAssetStreamingService() *AssetStreamingService {
	return &AssetStreamingService{}
}

func (s *AssetStreamingService) CanStream(asset *Asset, userID string, region string, userAge int) error {
	if asset == nil {
		return ErrInvalidAsset
	}

	if !asset.IsPublished() {
		if !asset.IsAccessibleBy(userID) {
			return ErrUnauthorizedAccess
		}
	}

	if !asset.IsPublished() {
		return ErrAssetNotPublished
	}

	if !asset.IsAvailableInRegion(region) {
		return ErrAssetNotAvailableInRegion
	}

	if !asset.IsAgeAppropriate(userAge) {
		return ErrAssetNotAgeAppropriate
	}

	if !asset.HasVideo() {
		return ErrAssetHasNoVideo
	}

	mainVideo := asset.GetMainVideo()
	if mainVideo == nil {
		return ErrAssetHasNoMainVideo
	}

	if !mainVideo.IsReady() {
		return ErrVideoNotReady
	}

	return nil
}

func (s *AssetStreamingService) GetStreamingInfo(asset *Asset, userID string, region string, userAge int) (*StreamingInfo, error) {
	if err := s.CanStream(asset, userID, region, userAge); err != nil {
		return nil, err
	}

	mainVideo := asset.GetMainVideo()
	if mainVideo == nil {
		return nil, ErrAssetHasNoMainVideo
	}

	thumbnail := asset.GetThumbnail()

	return &StreamingInfo{
		AssetID:      asset.ID().Value(),
		Title:        asset.Title().Value(),
		Description:  asset.Description().Value(),
		VideoID:      mainVideo.ID().Value(),
		VideoURL:     mainVideo.StorageLocation().URL(),
		ThumbnailURL: thumbnail.URL(),
		Duration:     mainVideo.Duration(),
		Width:        mainVideo.Width(),
		Height:       mainVideo.Height(),
		Format:       mainVideo.Format().Value(),
		StreamInfo:   mainVideo.StreamInfo(),
	}, nil
}

// TODO: This is going to be a big feature
func (s *AssetStreamingService) GetRecommendedAssets(asset *Asset, limit int) []*Asset {
	if asset == nil || limit <= 0 {
		return nil
	}

	recommendations := make([]*Asset, 0, limit)

	if asset.Genre() != nil {
		genre := asset.Genre().Value()
		recommendations = append(recommendations, s.getAssetsByGenre(genre, limit)...)
	}

	if asset.Genres() != nil {
		for _, genre := range asset.Genres().Values() {
			if len(recommendations) >= limit {
				break
			}
			recommendations = append(recommendations, s.getAssetsByGenre(genre.Value(), limit-len(recommendations))...)
		}
	}

	if len(recommendations) < limit && asset.Tags() != nil {
		for _, tag := range asset.Tags().Values() {
			if len(recommendations) >= limit {
				break
			}
			recommendations = append(recommendations, s.getAssetsByTag(tag, limit-len(recommendations))...)
		}
	}

	return recommendations[:min(len(recommendations), limit)]
}

func (s *AssetStreamingService) getAssetsByGenre(genre string, limit int) []*Asset {
	return nil
}

func (s *AssetStreamingService) getAssetsByTag(tag string, limit int) []*Asset {
	return nil
}

type AssetSearchService struct{}

func NewAssetSearchService() *AssetSearchService {
	return &AssetSearchService{}
}

// TODO: Implement this properly
func (s *AssetSearchService) SearchAssets(query string, filters *SearchFilters) ([]*Asset, error) {
	if query == "" && filters == nil {
		return nil, ErrInvalidSearchQuery
	}

	results := make([]*Asset, 0)

	if query != "" {
		results = append(results, s.searchByQuery(query)...)
	}

	if filters != nil {
		results = s.filterResults(results, filters)
	}

	return results, nil
}

func (s *AssetSearchService) searchByQuery(query string) []*Asset {
	return nil
}

func (s *AssetSearchService) filterResults(assets []*Asset, filters *SearchFilters) []*Asset {
	filtered := make([]*Asset, 0)

	for _, asset := range assets {
		if s.matchesFilters(asset, filters) {
			filtered = append(filtered, asset)
		}
	}

	return filtered
}

func (s *AssetSearchService) matchesFilters(asset *Asset, filters *SearchFilters) bool {
	if filters == nil {
		return true
	}

	if filters.AssetType != nil && !asset.Type().Equals(*filters.AssetType) {
		return false
	}

	if filters.Genre != nil && (asset.Genre() == nil || !asset.Genre().Equals(*filters.Genre)) {
		return false
	}

	if filters.Status != nil && (asset.Status() == nil || !asset.Status().Equals(*filters.Status)) {
		return false
	}

	if filters.OnlyPublic && !asset.IsPublished() {
		return false
	}

	if filters.OnlyPublished && !asset.IsPublished() {
		return false
	}

	if filters.OnlyReady && !asset.IsReady() {
		return false
	}

	if filters.HasVideo && !asset.HasVideo() {
		return false
	}

	if filters.HasImage && !asset.HasImage() {
		return false
	}

	return true
}

type StreamingInfo struct {
	AssetID      string
	Title        string
	Description  string
	VideoID      string
	VideoURL     string
	ThumbnailURL string
	Duration     *float64
	Width        *int
	Height       *int
	Format       string
	StreamInfo   *StreamInfoValue
}

type SearchFilters struct {
	AssetType     *AssetType
	Genre         *Genre
	Status        *Status
	OnlyPublic    bool
	OnlyPublished bool
	OnlyReady     bool
	HasVideo      bool
	HasImage      bool
}

var (
	ErrInvalidAsset              = pkgerrors.NewValidationError("invalid asset", nil)
	ErrUnauthorizedAccess        = pkgerrors.NewValidationError("unauthorized access", nil)
	ErrAssetNotReady             = pkgerrors.NewValidationError("asset not ready", nil)
	ErrAssetHasNoContent         = pkgerrors.NewValidationError("asset has no content", nil)
	ErrAssetNotPublished         = pkgerrors.NewValidationError("asset not published", nil)
	ErrAssetNotAvailableInRegion = pkgerrors.NewValidationError("asset not available in region", nil)
	ErrAssetNotAgeAppropriate    = pkgerrors.NewValidationError("asset not age appropriate", nil)
	ErrAssetHasNoVideo           = pkgerrors.NewValidationError("asset has no video", nil)
	ErrAssetHasNoMainVideo       = pkgerrors.NewValidationError("asset has no main video", nil)
	ErrVideoNotReady             = pkgerrors.NewValidationError("video not ready", nil)
	ErrInvalidSearchQuery        = pkgerrors.NewValidationError("invalid search query", nil)
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
