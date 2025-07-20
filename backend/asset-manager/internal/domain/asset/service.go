package asset

import (
	"context"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type AssetDomainService interface {
	ValidateAssetForPublishing(asset *Asset) error
	CalculateAssetStatus(asset *Asset) string
	ValidateVideoTranscoding(asset *Asset, videoID string) error
	ValidateAssetHierarchy(asset *Asset, parentID *AssetID) error
	ValidateAssetAccess(asset *Asset, userID string) error
	CalculateAssetMetrics(asset *Asset) AssetMetrics
	ValidateAssetMetadata(asset *Asset) error
	DetermineAssetVisibility(asset *Asset, userID string) bool
	ValidateAssetForStreaming(asset *Asset) error
	CalculateAssetStorageUsage(asset *Asset) StorageUsage
}

type AssetMetrics struct {
	TotalVideos      int
	ReadyVideos      int
	ProcessingVideos int
	FailedVideos     int
	TotalImages      int
	TotalCredits     int
	TotalDuration    float64
	TotalSize        int64
	StorageBuckets   map[string]int64
}

type StorageUsage struct {
	VideoStorage   int64
	ImageStorage   int64
	TotalStorage   int64
	StorageBuckets map[string]int64
}

type DomainService struct {
	repo AssetRepository
}

func NewDomainService(repo AssetRepository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

func (s *DomainService) ValidateAssetForPublishing(asset *Asset) error {
	if asset.Title() == nil {
		return errors.NewValidationError("asset title is required", nil)
	}

	if asset.Type() == nil {
		return errors.NewValidationError("asset type is required", nil)
	}

	if len(asset.Videos()) == 0 {
		return errors.NewValidationError("asset must have at least one video", nil)
	}

	hasReadyVideo := false
	for _, video := range asset.Videos() {
		if video.IsReady() {
			hasReadyVideo = true
			break
		}
	}

	if !hasReadyVideo {
		return errors.NewValidationError("asset must have at least one ready video", nil)
	}

	if asset.PublishRule() == nil {
		return errors.NewValidationError("publish rule is required", nil)
	}

	if asset.PublishRule().PublishAt() == nil {
		return errors.NewValidationError("publish date is required", nil)
	}

	return nil
}

func (s *DomainService) CalculateAssetStatus(asset *Asset) string {
	return asset.Status()
}

func (s *DomainService) ValidateVideoTranscoding(asset *Asset, videoID string) error {
	video, err := asset.GetVideo(videoID)
	if err != nil {
		return err
	}

	if video.IsReady() {
		return errors.NewValidationError("video is already ready", nil)
	}

	if video.IsFailed() {
		return errors.NewValidationError("video transcoding failed", nil)
	}

	return nil
}

func (s *DomainService) ValidateAssetHierarchy(asset *Asset, parentID *AssetID) error {
	if parentID == nil {
		return nil
	}

	if asset.ID().Equals(*parentID) {
		return errors.NewValidationError("asset cannot be its own parent", nil)
	}

	parent, err := s.repo.FindByID(context.Background(), parentID.Value())
	if err != nil {
		return errors.NewNotFoundError("parent asset not found", err)
	}

	if parent == nil {
		return errors.NewNotFoundError("parent asset not found", nil)
	}

	if parent.Status() != constants.AssetStatusPublished {
		return errors.NewValidationError("parent asset must be published", nil)
	}

	return nil
}

func (s *DomainService) ValidateAssetAccess(asset *Asset, userID string) error {
	if asset.OwnerID() == nil {
		return nil
	}

	if asset.OwnerID().Value() == userID {
		return nil
	}

	if asset.Status() == constants.AssetStatusPublished {
		return nil
	}

	return errors.NewForbiddenError("access denied", nil)
}

func (s *DomainService) CalculateAssetMetrics(asset *Asset) AssetMetrics {
	metrics := AssetMetrics{
		StorageBuckets: make(map[string]int64),
	}

	for _, video := range asset.Videos() {
		metrics.TotalVideos++
		metrics.TotalDuration += video.Duration()
		metrics.TotalSize += video.Size()

		switch video.Status() {
		case VideoStatus(constants.VideoStatusReady):
			metrics.ReadyVideos++
		case VideoStatus(constants.VideoStatusPending), VideoStatus(constants.VideoStatusAnalyzing), VideoStatus(constants.VideoStatusTranscoding):
			metrics.ProcessingVideos++
		case VideoStatus(constants.VideoStatusFailed):
			metrics.FailedVideos++
		}

		if video.StorageLocation().Bucket() != "" {
			metrics.StorageBuckets[video.StorageLocation().Bucket()] += video.Size()
		}
	}

	metrics.TotalImages = len(asset.Images())
	metrics.TotalCredits = len(asset.Credits())

	for _, image := range asset.Images() {
		if image.Size() != nil {
			metrics.TotalSize += *image.Size()
		}
		if image.StorageLocation() != nil && image.StorageLocation().Bucket() != "" {
			metrics.StorageBuckets[image.StorageLocation().Bucket()] += *image.Size()
		}
	}

	return metrics
}

func (s *DomainService) ValidateAssetMetadata(asset *Asset) error {
	if asset.Metadata() == nil {
		return nil
	}

	for key, value := range asset.Metadata() {
		if len(key) > 100 {
			return errors.NewValidationError("metadata key too long", nil)
		}

		if strValue, ok := value.(string); ok {
			if len(strValue) > 1000 {
				return errors.NewValidationError("metadata value too long", nil)
			}
		}
	}

	return nil
}

func (s *DomainService) DetermineAssetVisibility(asset *Asset, userID string) bool {
	if asset.Status() == constants.AssetStatusPublished {
		return true
	}

	if asset.OwnerID() != nil && asset.OwnerID().Value() == userID {
		return true
	}

	return false
}

func (s *DomainService) ValidateAssetForStreaming(asset *Asset) error {
	if asset.Status() != constants.AssetStatusPublished {
		return errors.NewValidationError("asset is not published", nil)
	}

	hasStreamableVideo := false
	for _, video := range asset.Videos() {
		if video.IsReady() && video.StreamInfo() != nil {
			hasStreamableVideo = true
			break
		}
	}

	if !hasStreamableVideo {
		return errors.NewValidationError("asset has no streamable videos", nil)
	}

	return nil
}

func (s *DomainService) CalculateAssetStorageUsage(asset *Asset) StorageUsage {
	usage := StorageUsage{
		StorageBuckets: make(map[string]int64),
	}

	for _, video := range asset.Videos() {
		usage.VideoStorage += video.Size()
		usage.TotalStorage += video.Size()

		if video.StorageLocation().Bucket() != "" {
			usage.StorageBuckets[video.StorageLocation().Bucket()] += video.Size()
		}
	}

	for _, image := range asset.Images() {
		if image.Size() != nil {
			usage.ImageStorage += *image.Size()
			usage.TotalStorage += *image.Size()
		}

		if image.StorageLocation() != nil && image.StorageLocation().Bucket() != "" {
			usage.StorageBuckets[image.StorageLocation().Bucket()] += *image.Size()
		}
	}

	return usage
}

type AssetPublishingService interface {
	ValidatePublishingRules(asset *Asset) error
	CalculatePublishDate(publishRule *PublishRule) *time.Time
	ValidateAgeRating(asset *Asset, userAge int) error
	ValidateRegionalAccess(asset *Asset, userRegion string) error
}

type PublishingService struct {
	domainService AssetDomainService
}

func NewPublishingService(domainService AssetDomainService) *PublishingService {
	return &PublishingService{
		domainService: domainService,
	}
}

func (s *PublishingService) ValidatePublishingRules(asset *Asset) error {
	if asset.PublishRule() == nil {
		return errors.NewValidationError("publish rule is required", nil)
	}

	if asset.PublishRule().PublishAt() == nil {
		return errors.NewValidationError("publish date is required", nil)
	}

	if asset.PublishRule().UnpublishAt() != nil {
		if asset.PublishRule().PublishAt().After(*asset.PublishRule().UnpublishAt()) {
			return errors.NewValidationError("publish date must be before unpublish date", nil)
		}
	}

	return nil
}

func (s *PublishingService) CalculatePublishDate(publishRule *PublishRule) *time.Time {
	if publishRule == nil || publishRule.PublishAt() == nil {
		return nil
	}

	now := time.Now().UTC()
	if now.Before(*publishRule.PublishAt()) {
		return publishRule.PublishAt()
	}

	return &now
}

func (s *PublishingService) ValidateAgeRating(asset *Asset, userAge int) error {
	if asset.PublishRule() == nil || asset.PublishRule().AgeRating() == nil {
		return nil
	}

	ageRating := *asset.PublishRule().AgeRating()
	requiredAge := s.getRequiredAge(ageRating)

	if userAge < requiredAge {
		return errors.NewForbiddenError("age restriction", nil)
	}

	return nil
}

func (s *PublishingService) ValidateRegionalAccess(asset *Asset, userRegion string) error {
	if asset.PublishRule() == nil || len(asset.PublishRule().Regions()) == 0 {
		return nil
	}

	for _, region := range asset.PublishRule().Regions() {
		if region == userRegion {
			return nil
		}
	}

	return errors.NewForbiddenError("regional restriction", nil)
}

func (s *PublishingService) getRequiredAge(ageRating string) int {
	ageMap := map[string]int{
		constants.AgeRatingG: 0, constants.AgeRatingPG: 0, constants.AgeRatingPG13: 13, constants.AgeRatingR: 17, constants.AgeRatingNC17: 18,
		constants.AgeRatingTVY: 0, constants.AgeRatingTVY7: 7, constants.AgeRatingTVG: 0, constants.AgeRatingTVPG: 0, constants.AgeRatingTV14: 14, constants.AgeRatingTVMA: 17,
	}

	if age, exists := ageMap[ageRating]; exists {
		return age
	}

	return 0
}
