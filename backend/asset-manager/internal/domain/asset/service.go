package asset

import (
	"context"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type AssetDomainService interface {
	ValidateAssetHierarchy(asset *Asset, parentID *AssetID) error
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

func (s *DomainService) ValidateAssetHierarchy(asset *Asset, parentID *AssetID) error {
	if parentID == nil {
		return nil
	}

	if asset.ID().Equals(*parentID) {
		return errors.NewValidationError("asset cannot be its own parent", nil)
	}

	parent, err := s.repo.FindByID(context.Background(), *parentID)
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
