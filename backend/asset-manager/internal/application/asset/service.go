package asset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	domainasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

type AssetApplicationService interface {
	CreateAsset(ctx context.Context, cmd CreateAssetCommand) (*domainasset.Asset, error)
	PatchAsset(ctx context.Context, cmd PatchAssetCommand) (*domainasset.Asset, error)
	UpdateAsset(ctx context.Context, cmd UpdateAssetCommand) error
	DeleteAsset(ctx context.Context, cmd DeleteAssetCommand) error
	GetAsset(ctx context.Context, query GetAssetQuery) (*domainasset.Asset, error)
	ListAssets(ctx context.Context, query ListAssetsQuery) (*domainasset.AssetPage, error)
	SearchAssets(ctx context.Context, query SearchAssetsQuery) (*domainasset.AssetPage, error)
	AddVideo(ctx context.Context, cmd AddVideoCommand) (*domainasset.Video, error)
	RemoveVideo(ctx context.Context, cmd RemoveVideoCommand) error
	UpdateVideoStatus(ctx context.Context, cmd UpdateVideoStatusCommand) error
	AddImage(ctx context.Context, cmd AddImageCommand) error
	RemoveImage(ctx context.Context, cmd RemoveImageCommand) error
	PublishAsset(ctx context.Context, cmd PublishAssetCommand) error
	UpdateVideoAnalysis(ctx context.Context, assetID, videoID string, metadata *messages.JobCompletionPayload) error
	UpdateVideoTranscoding(ctx context.Context, assetID, videoID, format string, metadata *messages.JobCompletionPayload) error
	GetAssetMetrics(ctx context.Context, assetID string) (*domainasset.AssetMetrics, error)
	GetAssetStorageUsage(ctx context.Context, assetID string) (*domainasset.StorageUsage, error)
	ValidateAssetAccess(ctx context.Context, assetID, userID string) error
}

type ApplicationService struct {
	repo              domainasset.AssetRepository
	domainService     domainasset.AssetDomainService
	publishingService domainasset.AssetPublishingService
	eventPublisher    EventPublisher
	logger            *logger.Logger
}

func NewApplicationService(
	repo domainasset.AssetRepository,
	domainService domainasset.AssetDomainService,
	publishingService domainasset.AssetPublishingService,
	eventPublisher EventPublisher,
) *ApplicationService {
	return &ApplicationService{
		repo:              repo,
		domainService:     domainService,
		publishingService: publishingService,
		eventPublisher:    eventPublisher,
		logger:            logger.WithService("asset-application-service"),
	}
}

func (s *ApplicationService) CreateAsset(ctx context.Context, cmd CreateAssetCommand) (*domainasset.Asset, error) {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid create asset command")
		return nil, pkgerrors.NewValidationError("invalid command", err)
	}

	slug, title, description, assetType, genre, genres, tags, ownerID, parentID, err := cmd.ToDomainValueObjects()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain value objects")
		return nil, pkgerrors.NewValidationError("invalid value objects", err)
	}

	asset, err := domainasset.NewAsset(*slug, title, assetType)
	if err != nil {
		log.WithError(err).Error("Failed to create asset", "slug", slug.Value())
		return nil, pkgerrors.NewValidationError("failed to create asset", err)
	}

	if description != nil {
		if err := asset.UpdateDescription(description); err != nil {
			log.WithError(err).Error("Failed to update description")
			return nil, err
		}
	}

	if genre != nil {
		if err := asset.UpdateGenre(genre); err != nil {
			log.WithError(err).Error("Failed to update genre")
			return nil, err
		}
	}

	if genres != nil {
		if err := asset.SetGenres(genres); err != nil {
			log.WithError(err).Error("Failed to set genres")
			return nil, err
		}
	}

	if tags != nil {
		if err := asset.SetTags(tags); err != nil {
			log.WithError(err).Error("Failed to set tags")
			return nil, err
		}
	}

	if ownerID != nil {
		if err := asset.SetOwnerID(ownerID); err != nil {
			log.WithError(err).Error("Failed to set owner ID")
			return nil, err
		}
	}

	if parentID != nil {
		if err := s.domainService.ValidateAssetHierarchy(asset, parentID); err != nil {
			log.WithError(err).Error("Failed to validate asset hierarchy", "parent_id", parentID.Value())
			return nil, err
		}
		if err := asset.SetParentID(parentID); err != nil {
			log.WithError(err).Error("Failed to set parent ID")
			return nil, err
		}
	}

	if cmd.PublishRule != nil {
		if err := asset.SetPublishRule(cmd.PublishRule); err != nil {
			log.WithError(err).Error("Failed to set publish rule")
			return nil, err
		}
	}

	if cmd.Metadata != nil {
		if err := asset.SetMetadata(cmd.Metadata); err != nil {
			log.WithError(err).Error("Failed to set metadata")
			return nil, err
		}
	}

	if err := s.domainService.ValidateAssetMetadata(asset); err != nil {
		log.WithError(err).Error("Failed to validate asset metadata")
		return nil, err
	}

	if err := s.repo.Save(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to save asset", "asset_id", asset.ID().Value())
		return nil, err
	}

	if err := s.eventPublisher.PublishAssetCreated(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to publish asset created event", "asset_id", asset.ID().Value())
	}
	log.Info("Asset created successfully", "asset_id", asset.ID().Value(), "slug", asset.Slug().Value())
	return asset, nil
}

func (s *ApplicationService) UpdateAsset(ctx context.Context, cmd UpdateAssetCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid update asset command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, title, description, assetType, genre, genres, tags, ownerID, parentID, err := cmd.ToDomainValueObjects()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain value objects")
		return pkgerrors.NewValidationError("invalid value objects", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for update", "asset_id", assetID.Value())
		return err
	}

	if title != nil {
		if err := asset.UpdateTitle(title); err != nil {
			log.WithError(err).Error("Failed to update title")
			return err
		}
	}

	if description != nil {
		if err := asset.UpdateDescription(description); err != nil {
			log.WithError(err).Error("Failed to update description")
			return err
		}
	}

	if assetType != nil {
		if err := asset.UpdateType(assetType); err != nil {
			log.WithError(err).Error("Failed to update type")
			return err
		}
	}

	if genre != nil {
		if err := asset.UpdateGenre(genre); err != nil {
			log.WithError(err).Error("Failed to update genre")
			return err
		}
	}

	if genres != nil {
		if err := asset.SetGenres(genres); err != nil {
			log.WithError(err).Error("Failed to set genres")
			return err
		}
	}

	if tags != nil {
		if err := asset.SetTags(tags); err != nil {
			log.WithError(err).Error("Failed to set tags")
			return err
		}
	}

	if ownerID != nil {
		if err := asset.SetOwnerID(ownerID); err != nil {
			log.WithError(err).Error("Failed to set owner ID")
			return err
		}
	}

	if parentID != nil {
		if err := s.domainService.ValidateAssetHierarchy(asset, parentID); err != nil {
			log.WithError(err).Error("Failed to validate asset hierarchy", "parent_id", parentID.Value())
			return err
		}
		if err := asset.SetParentID(parentID); err != nil {
			log.WithError(err).Error("Failed to set parent ID")
			return err
		}
	}

	if cmd.PublishRule != nil {
		if err := asset.SetPublishRule(cmd.PublishRule); err != nil {
			log.WithError(err).Error("Failed to set publish rule")
			return err
		}
	}

	if cmd.Metadata != nil {
		if err := asset.SetMetadata(cmd.Metadata); err != nil {
			log.WithError(err).Error("Failed to set metadata")
			return err
		}
	}

	if err := s.domainService.ValidateAssetMetadata(asset); err != nil {
		log.WithError(err).Error("Failed to validate asset metadata")
		return err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishAssetUpdated(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to publish asset updated event", "asset_id", assetID.Value())
	}
	log.Info("Asset updated successfully", "asset_id", assetID.Value())
	return nil
}

func (s *ApplicationService) PatchAsset(ctx context.Context, cmd PatchAssetCommand) (*domainasset.Asset, error) {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid patch asset command")
		return nil, pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return nil, pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for patching", "asset_id", assetID.Value())
		return nil, err
	}

	for _, patch := range cmd.Patches {
		log.Info("Applying patch", "patch", patch, "asset_id", assetID.Value(), "asset_status", asset.Status())
		if err := s.applyPatch(asset, patch); err != nil {
			log.WithError(err).Error("Failed to apply patch", "patch", patch, "asset_id", assetID.Value())
			return nil, pkgerrors.NewValidationError("failed to apply patch", err)
		}
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after patching", "asset_id", assetID.Value())
		return nil, err
	}

	log.Info("Asset patched successfully", "asset_id", assetID.Value())
	return asset, nil
}

func (s *ApplicationService) applyPatch(a *domainasset.Asset, patch JSONPatchOperation) error {
	switch patch.Op {
	case "replace":
		return s.applyReplacePatch(a, patch)
	case "remove":
		return s.applyRemovePatch(a, patch)
	default:
		return errors.New("unsupported patch operation")
	}
}

func (s *ApplicationService) applyReplacePatch(a *domainasset.Asset, patch JSONPatchOperation) error {
	patchHandlers := map[string]func(*domainasset.Asset, string) error{
		"/title":       s.patchTitle,
		"/description": s.patchDescription,
		"/type":        s.patchType,
		"/genre":       s.patchGenre,
		"/ownerId":     s.patchOwnerID,
		"/publishAt":   s.patchPublishAt,
		"/unpublishAt": s.patchUnpublishAt,
		"/regions":     s.patchRegions,
		"/ageRating":   s.patchAgeRating,
	}

	handler, exists := patchHandlers[patch.Path]
	if !exists {
		return errors.New("unsupported field for replacement")
	}

	return handler(a, patch.Value)
}

func (s *ApplicationService) applyRemovePatch(a *domainasset.Asset, patch JSONPatchOperation) error {
	removeHandlers := map[string]func(*domainasset.Asset) error{
		"/title":       s.removeTitle,
		"/description": s.removeDescription,
		"/genre":       s.removeGenre,
		"/ownerId":     s.removeOwnerID,
		"/publishAt":   s.removePublishAt,
		"/unpublishAt": s.removeUnpublishAt,
		"/regions":     s.removeRegions,
		"/ageRating":   s.removeAgeRating,
	}

	handler, exists := removeHandlers[patch.Path]
	if !exists {
		return errors.New("unsupported field for removal")
	}

	return handler(a)
}

func (s *ApplicationService) patchTitle(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.UpdateTitle(nil)
	}
	title, err := domainasset.NewTitle(value)
	if err != nil {
		return err
	}
	return a.UpdateTitle(title)
}

func (s *ApplicationService) patchDescription(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.UpdateDescription(nil)
	}
	description, err := domainasset.NewDescription(value)
	if err != nil {
		return err
	}
	return a.UpdateDescription(description)
}

func (s *ApplicationService) patchType(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.UpdateType(nil)
	}
	assetType, err := domainasset.NewAssetType(value)
	if err != nil {
		return err
	}
	return a.UpdateType(assetType)
}

func (s *ApplicationService) patchGenre(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.UpdateGenre(nil)
	}
	genre, err := domainasset.NewGenre(value)
	if err != nil {
		return err
	}
	return a.UpdateGenre(genre)
}

func (s *ApplicationService) patchOwnerID(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.SetOwnerID(nil)
	}
	ownerID, err := domainasset.NewOwnerID(value)
	if err != nil {
		return err
	}
	return a.SetOwnerID(ownerID)
}

func (s *ApplicationService) patchPublishAt(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.SetPublishRule(nil)
	}
	publishTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	return s.updatePublishRule(a, &publishTime, nil, nil, nil)
}

func (s *ApplicationService) patchUnpublishAt(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.SetPublishRule(nil)
	}
	unpublishTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	return s.updatePublishRule(a, nil, &unpublishTime, nil, nil)
}

func (s *ApplicationService) patchRegions(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.SetPublishRule(nil)
	}
	var regions []string
	if err := json.Unmarshal([]byte(value), &regions); err != nil {
		return err
	}
	return s.updatePublishRule(a, nil, nil, regions, nil)
}

func (s *ApplicationService) patchAgeRating(a *domainasset.Asset, value string) error {
	if value == "" {
		return a.SetPublishRule(nil)
	}
	return s.updatePublishRule(a, nil, nil, nil, &value)
}

func (s *ApplicationService) removeTitle(a *domainasset.Asset) error {
	return a.UpdateTitle(nil)
}

func (s *ApplicationService) removeDescription(a *domainasset.Asset) error {
	return a.UpdateDescription(nil)
}

func (s *ApplicationService) removeGenre(a *domainasset.Asset) error {
	return a.UpdateGenre(nil)
}

func (s *ApplicationService) removeOwnerID(a *domainasset.Asset) error {
	return a.SetOwnerID(nil)
}

func (s *ApplicationService) removePublishAt(a *domainasset.Asset) error {
	return s.updatePublishRule(a, nil, nil, nil, nil)
}

func (s *ApplicationService) removeUnpublishAt(a *domainasset.Asset) error {
	return s.updatePublishRule(a, nil, nil, nil, nil)
}

func (s *ApplicationService) removeRegions(a *domainasset.Asset) error {
	return s.updatePublishRule(a, nil, nil, []string{}, nil)
}

func (s *ApplicationService) removeAgeRating(a *domainasset.Asset) error {
	return s.updatePublishRule(a, nil, nil, nil, nil)
}

func (s *ApplicationService) updatePublishRule(a *domainasset.Asset, publishAt, unpublishAt *time.Time, regions []string, ageRating *string) error {
	var existingRule *domainasset.PublishRule
	if a.PublishRule() != nil {
		existingRule = a.PublishRule()
	}

	// Use provided values or fall back to existing ones
	if publishAt == nil && existingRule != nil {
		publishAt = existingRule.PublishAt()
	}
	if unpublishAt == nil && existingRule != nil {
		unpublishAt = existingRule.UnpublishAt()
	}
	if regions == nil && existingRule != nil {
		regions = existingRule.Regions()
	}
	if ageRating == nil && existingRule != nil {
		ageRating = existingRule.AgeRating()
	}

	publishRule, err := domainasset.NewPublishRule(publishAt, unpublishAt, regions, ageRating)
	if err != nil {
		return err
	}
	return a.SetPublishRule(publishRule)
}

func (s *ApplicationService) DeleteAsset(ctx context.Context, cmd DeleteAssetCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid delete asset command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for deletion", "asset_id", assetID.Value())
		return err
	}

	if err := s.repo.Delete(ctx, assetID.Value()); err != nil {
		log.WithError(err).Error("Failed to delete asset", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishAssetDeleted(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to publish asset deleted event", "asset_id", assetID.Value())
	}
	log.Info("Asset deleted successfully", "asset_id", assetID.Value())
	return nil
}

func (s *ApplicationService) GetAsset(ctx context.Context, query GetAssetQuery) (*domainasset.Asset, error) {
	log := s.logger.WithContext(ctx)

	if err := query.Validate(); err != nil {
		log.WithError(err).Error("Invalid get asset query")
		return nil, pkgerrors.NewValidationError("invalid query", err)
	}

	var asset *domainasset.Asset
	var err error

	if query.ID != "" {
		assetID, err := query.ToDomainAssetID()
		if err != nil {
			log.WithError(err).Error("Failed to convert query to domain asset ID")
			return nil, pkgerrors.NewValidationError("invalid asset ID", err)
		}
		asset, err = s.repo.FindByID(ctx, assetID.Value())
	} else {
		slug, err := query.ToDomainSlug()
		if err != nil {
			log.WithError(err).Error("Failed to convert query to domain slug")
			return nil, pkgerrors.NewValidationError("invalid slug", err)
		}
		asset, err = s.repo.FindBySlug(ctx, slug.Value())
	}

	if err != nil {
		log.WithError(err).Error("Failed to get asset", "id", query.ID, "slug", query.Slug)
		return nil, err
	}

	return asset, nil
}

func (s *ApplicationService) ListAssets(ctx context.Context, query ListAssetsQuery) (*domainasset.AssetPage, error) {
	log := s.logger.WithContext(ctx)

	if err := query.Validate(); err != nil {
		log.WithError(err).Error("Invalid list assets query")
		return nil, pkgerrors.NewValidationError("invalid query", err)
	}

	assets, err := s.repo.List(ctx, query.Limit, query.LastKey)
	if err != nil {
		log.WithError(err).Error("Failed to list assets")
		return nil, err
	}

	return assets, nil
}

func (s *ApplicationService) SearchAssets(ctx context.Context, query SearchAssetsQuery) (*domainasset.AssetPage, error) {
	log := s.logger.WithContext(ctx)

	if err := query.Validate(); err != nil {
		log.WithError(err).Error("Invalid search assets query")
		return nil, pkgerrors.NewValidationError("invalid query", err)
	}

	assets, err := s.repo.Search(ctx, query.Query, query.Limit, query.LastKey)
	if err != nil {
		log.WithError(err).Error("Failed to search assets")
		return nil, err
	}

	return assets, nil
}

func (s *ApplicationService) AddVideo(ctx context.Context, cmd AddVideoCommand) (*domainasset.Video, error) {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid add video command")
		return nil, pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return nil, pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for video addition", "asset_id", assetID.Value())
		return nil, err
	}

	format, err := domainasset.NewVideoFormat(cmd.Format)
	if err != nil {
		log.WithError(err).Error("Failed to create video format", "asset_id", assetID.Value(), "format", cmd.Format)
		return nil, pkgerrors.NewValidationError("invalid video format", err)
	}

	video, err := asset.AddVideo(cmd.Label, format, cmd.StorageLocation)
	if err != nil {
		log.WithError(err).Error("Failed to add video to asset", "asset_id", assetID.Value(), "label", cmd.Label)
		return nil, err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after adding video", "asset_id", assetID.Value())
		return nil, err
	}

	if err := s.eventPublisher.PublishVideoAdded(ctx, asset, video); err != nil {
		log.WithError(err).Error("Failed to publish video added event", "asset_id", assetID.Value(), "video_id", video.ID())
	}
	log.Info("Video added successfully", "asset_id", assetID.Value(), "video_id", video.ID())
	return video, nil
}

func (s *ApplicationService) RemoveVideo(ctx context.Context, cmd RemoveVideoCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid remove video command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for video removal", "asset_id", assetID.Value())
		return err
	}

	if err := asset.RemoveVideo(cmd.VideoID); err != nil {
		log.WithError(err).Error("Failed to remove video from asset", "asset_id", assetID.Value(), "video_id", cmd.VideoID)
		return err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after removing video", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishVideoRemoved(ctx, asset, cmd.VideoID); err != nil {
		log.WithError(err).Error("Failed to publish video removed event", "asset_id", assetID.Value(), "video_id", cmd.VideoID)
	}
	log.Info("Video removed successfully", "asset_id", assetID.Value(), "video_id", cmd.VideoID)
	return nil
}

func (s *ApplicationService) UpdateVideoStatus(ctx context.Context, cmd UpdateVideoStatusCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid update video status command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for video status update", "asset_id", assetID.Value())
		return err
	}

	if err := asset.UpdateVideoStatus(cmd.VideoID, cmd.Status); err != nil {
		log.WithError(err).Error("Failed to update video status", "asset_id", assetID.Value(), "video_id", cmd.VideoID)
		return err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after video status change", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishVideoStatusUpdated(ctx, asset, cmd.VideoID, cmd.Status); err != nil {
		log.WithError(err).Error("Failed to publish video status updated event", "asset_id", assetID.Value(), "video_id", cmd.VideoID)
	}
	log.Info("Video status updated successfully", "asset_id", assetID.Value(), "video_id", cmd.VideoID, "status", cmd.Status)
	return nil
}

func (s *ApplicationService) AddImage(ctx context.Context, cmd AddImageCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid add image command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for image addition", "asset_id", assetID.Value())
		return err
	}

	if err := asset.AddImage(cmd.Image); err != nil {
		log.WithError(err).Error("Failed to add image to asset", "asset_id", assetID.Value(), "image_id", cmd.Image.ID())
		return err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after adding image", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishImageAdded(ctx, asset, cmd.Image); err != nil {
		log.WithError(err).Error("Failed to publish image added event", "asset_id", assetID.Value(), "image_id", cmd.Image.ID())
	}
	log.Info("Image added successfully", "asset_id", assetID.Value(), "image_id", cmd.Image.ID())
	return nil
}

func (s *ApplicationService) RemoveImage(ctx context.Context, cmd RemoveImageCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid remove image command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for image removal", "asset_id", assetID.Value())
		return err
	}

	if err := asset.RemoveImage(cmd.ImageID); err != nil {
		log.WithError(err).Error("Failed to remove image from asset", "asset_id", assetID.Value(), "image_id", cmd.ImageID)
		return err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after removing image", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishImageRemoved(ctx, asset, cmd.ImageID); err != nil {
		log.WithError(err).Error("Failed to publish image removed event", "asset_id", assetID.Value(), "image_id", cmd.ImageID)
	}
	log.Info("Image removed successfully", "asset_id", assetID.Value(), "image_id", cmd.ImageID)
	return nil
}

func (s *ApplicationService) PublishAsset(ctx context.Context, cmd PublishAssetCommand) error {
	log := s.logger.WithContext(ctx)

	if err := cmd.Validate(); err != nil {
		log.WithError(err).Error("Invalid publish asset command")
		return pkgerrors.NewValidationError("invalid command", err)
	}

	assetID, err := cmd.ToDomainAssetID()
	if err != nil {
		log.WithError(err).Error("Failed to convert command to domain asset ID")
		return pkgerrors.NewValidationError("invalid asset ID", err)
	}

	asset, err := s.repo.FindByID(ctx, assetID.Value())
	if err != nil {
		log.WithError(err).Error("Failed to find asset for publishing", "asset_id", assetID.Value())
		return err
	}

	if err := s.domainService.ValidateAssetForPublishing(asset); err != nil {
		log.WithError(err).Error("Failed to validate asset for publishing", "asset_id", assetID.Value())
		return err
	}

	if err := s.publishingService.ValidatePublishingRules(asset); err != nil {
		log.WithError(err).Error("Failed to validate publishing rules", "asset_id", assetID.Value())
		return err
	}

	if err := asset.SetPublishRule(cmd.PublishRule); err != nil {
		log.WithError(err).Error("Failed to set publish rule", "asset_id", assetID.Value())
		return err
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset for publishing", "asset_id", assetID.Value())
		return err
	}

	if err := s.eventPublisher.PublishAssetPublished(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to publish asset published event", "asset_id", assetID.Value())
	}
	log.Info("Asset published successfully", "asset_id", assetID.Value())
	return nil
}

func (s *ApplicationService) UpdateVideoAnalysis(ctx context.Context, assetID, videoID string, metadata *messages.JobCompletionPayload) error {
	log := s.logger.WithContext(ctx)

	asset, err := s.repo.FindByID(ctx, assetID)
	if err != nil {
		log.WithError(err).Error("Failed to find asset for video analysis update", "asset_id", assetID)
		return err
	}

	video, err := asset.GetVideo(videoID)
	if err != nil {
		log.WithError(err).Error("Failed to get video for analysis update", "asset_id", assetID, "video_id", videoID)
		return err
	}

	if metadata.Success {
		video.UpdateStatus(domainasset.VideoStatus(constants.VideoStatusReady))
		if metadata.Duration > 0 {
			video.UpdateDuration(metadata.Duration)
		}
		if metadata.Width > 0 && metadata.Height > 0 {
			video.UpdateDimensions(metadata.Width, metadata.Height)
		}
		if metadata.Bitrate > 0 {
			video.UpdateBitrate(metadata.Bitrate)
		}
		if metadata.Codec != "" {
			video.UpdateCodec(metadata.Codec)
		}
		if metadata.Size > 0 {
			video.UpdateSize(metadata.Size)
		}
		if metadata.ContentType != "" {
			video.UpdateContentType(metadata.ContentType)
		}
	} else {
		video.UpdateStatus(domainasset.VideoStatus(constants.VideoStatusFailed))
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after video analysis", "asset_id", assetID)
		return err
	}

	if err := s.eventPublisher.PublishVideoStatusUpdated(ctx, asset, videoID, video.Status()); err != nil {
		log.WithError(err).Error("Failed to publish video status updated event", "asset_id", assetID, "video_id", videoID)
	}
	log.Info("Video analysis updated successfully", "asset_id", assetID, "video_id", videoID, "success", metadata.Success)
	return nil
}

func (s *ApplicationService) UpdateVideoTranscoding(ctx context.Context, assetID, videoID, format string, metadata *messages.JobCompletionPayload) error {
	log := s.logger.WithContext(ctx)

	asset, err := s.repo.FindByID(ctx, assetID)
	if err != nil {
		log.WithError(err).Error("Failed to find asset for video transcoding update", "asset_id", assetID)
		return err
	}

	video, err := asset.GetVideo(videoID)
	if err != nil {
		log.WithError(err).Error("Failed to get video for transcoding update", "asset_id", assetID, "video_id", videoID)
		return err
	}

	if metadata.Success {
		video.UpdateStatus(domainasset.VideoStatus(constants.VideoStatusReady))

		if format == "hls" || format == "dash" {
			if metadata.Bucket != "" && metadata.Key != "" {
				manifestS3Url := "s3://" + metadata.Bucket + "/" + metadata.Key
				newS3Obj, err := domainasset.NewS3Object(metadata.Bucket, metadata.Key, manifestS3Url)
				if err == nil {
					video.SetStorageLocation(*newS3Obj)
				}
			}
			if metadata.URL != "" {
				cdnURL := s.convertS3URLToCDN(metadata.URL)
				cdnPrefix := "http://localhost:8083/cdn"
				log.Info("Creating StreamInfo", "metadata_url", metadata.URL, "cdn_url", cdnURL, "cdn_prefix", cdnPrefix)
				streamInfo, err := domainasset.NewStreamInfo(&metadata.URL, &cdnPrefix, &cdnURL)
				if err != nil {
					log.WithError(err).Error("Failed to create StreamInfo", "metadata_url", metadata.URL, "cdn_url", cdnURL)
				} else {
					log.Info("Successfully created StreamInfo", "stream_info", streamInfo)
					video.SetStreamInfo(streamInfo)
				}
			}
		}
	} else {
		video.UpdateStatus(domainasset.VideoStatus(constants.VideoStatusFailed))
	}

	transcodingInfo, err := domainasset.NewTranscodingInfo(
		metadata.VideoID,
		100.0,
		metadata.URL,
		&metadata.Error,
		nil,
	)
	if err == nil {
		video.UpdateTranscodingInfo(*transcodingInfo)
	}

	if err := s.repo.Update(ctx, asset); err != nil {
		log.WithError(err).Error("Failed to update asset after video transcoding", "asset_id", assetID)
		return err
	}

	if err := s.eventPublisher.PublishVideoStatusUpdated(ctx, asset, videoID, video.Status()); err != nil {
		log.WithError(err).Error("Failed to publish video status updated event", "asset_id", assetID, "video_id", videoID)
	}
	log.Info("Video transcoding updated successfully", "asset_id", assetID, "video_id", videoID, "format", format, "success", metadata.Success)
	return nil
}

func (s *ApplicationService) convertS3URLToCDN(s3URL string) string {
	if !strings.HasPrefix(s3URL, "s3://") {
		return s3URL
	}

	s3Path := strings.TrimPrefix(s3URL, "s3://")
	parts := strings.SplitN(s3Path, "/", 2)
	if len(parts) != 2 {
		return s3URL
	}

	key := parts[1]

	return fmt.Sprintf("http://localhost:8083/cdn/%s", key)
}

func (s *ApplicationService) GetAssetMetrics(ctx context.Context, assetID string) (*domainasset.AssetMetrics, error) {
	log := s.logger.WithContext(ctx)

	asset, err := s.repo.FindByID(ctx, assetID)
	if err != nil {
		log.WithError(err).Error("Failed to find asset for metrics", "asset_id", assetID)
		return nil, err
	}

	metrics := s.domainService.CalculateAssetMetrics(asset)
	return &metrics, nil
}

func (s *ApplicationService) GetAssetStorageUsage(ctx context.Context, assetID string) (*domainasset.StorageUsage, error) {
	log := s.logger.WithContext(ctx)

	asset, err := s.repo.FindByID(ctx, assetID)
	if err != nil {
		log.WithError(err).Error("Failed to find asset for storage usage", "asset_id", assetID)
		return nil, err
	}

	usage := s.domainService.CalculateAssetStorageUsage(asset)
	return &usage, nil
}

func (s *ApplicationService) ValidateAssetAccess(ctx context.Context, assetID, userID string) error {
	log := s.logger.WithContext(ctx)

	asset, err := s.repo.FindByID(ctx, assetID)
	if err != nil {
		log.WithError(err).Error("Failed to find asset for access validation", "asset_id", assetID)
		return err
	}

	return s.domainService.ValidateAssetAccess(asset, userID)
}
