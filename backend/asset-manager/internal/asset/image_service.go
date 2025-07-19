package asset

import (
	"context"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type ImageService struct {
	repo   AssetRepository
	config *config.DynamicConfig
}

func NewImageService(repo AssetRepository, config *config.DynamicConfig) *ImageService {
	return &ImageService{
		repo:   repo,
		config: config,
	}
}

func (s *ImageService) AddImage(ctx context.Context, assetID string, img *Image) error {
	asset, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	if img.ID == "" {
		img.ID = generateID()
	}

	img.CreatedAt = time.Now()
	img.UpdatedAt = time.Now()

	for _, existing := range asset.Images {
		if existing.FileName == img.FileName && existing.Type == img.Type {
			return apperrors.NewConflictError("image with same filename and type already exists for this asset", nil)
		}
	}

	if img.StorageLocation != nil {
		cdnPrefix := s.getCDNPrefixForBucket(img.StorageLocation.Bucket)
		if cdnPrefix != "" {
			url := cdnPrefix + "/" + img.StorageLocation.Key
			img.StreamInfo = &StreamInfo{
				CdnPrefix: &cdnPrefix,
				URL:       &url,
			}
		}
	}

	asset.Images = append(asset.Images, *img)
	return s.repo.SaveAsset(ctx, asset)
}

func (s *ImageService) DeleteImage(ctx context.Context, assetID string, imageID string) error {
	asset, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return apperrors.NewNotFoundError("asset not found", err)
	}

	filtered := make([]Image, 0, len(asset.Images))
	for _, img := range asset.Images {
		if img.ID != imageID {
			filtered = append(filtered, img)
		}
	}

	asset.Images = filtered
	return s.repo.SaveAsset(ctx, asset)
}

func (s *ImageService) getCDNPrefixForBucket(bucket string) string {
	switch bucket {
	case "content-east", "content-west":
		return s.config.GetStringFromComponent("cdn", "prefix")
	default:
		return ""
	}
}
