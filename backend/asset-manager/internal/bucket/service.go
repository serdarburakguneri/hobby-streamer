package bucket

import (
	"context"
	"time"

	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type BucketService interface {
	GetBucketByID(ctx context.Context, id string) (*Bucket, error)
	GetBucketByKey(ctx context.Context, key string) (*Bucket, error)
	ListBuckets(ctx context.Context, limit int) (*BucketPage, error)
	CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error)
	UpdateBucket(ctx context.Context, id string, b *Bucket) error
	PatchBucket(ctx context.Context, id string, patch map[string]interface{}) error
	DeleteBucket(ctx context.Context, id string) error
	AddAssetToBucket(ctx context.Context, bucketID string, assetID string) error
	RemoveAssetFromBucket(ctx context.Context, bucketID string, assetID string) error
	GetBucketsByType(ctx context.Context, bucketType string) ([]Bucket, error)
	GetBucketsByAsset(ctx context.Context, assetID string) ([]Bucket, error)
	GetAssetsInBucket(ctx context.Context, bucketID string) ([]string, error)
}

type Service struct {
	Repo BucketRepository
}

var _ BucketService = (*Service)(nil)

func (s *Service) GetBucketByID(ctx context.Context, id string) (*Bucket, error) {
	bucket, err := s.Repo.GetBucketByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("bucket not found", err)
	}
	return bucket, nil
}

func (s *Service) GetBucketByKey(ctx context.Context, key string) (*Bucket, error) {
	bucket, err := s.Repo.GetBucketByKey(ctx, key)
	if err != nil {
		return nil, apperrors.NewNotFoundError("bucket not found", err)
	}
	return bucket, nil
}

func (s *Service) ListBuckets(ctx context.Context, limit int) (*BucketPage, error) {
	return s.Repo.ListBuckets(ctx, limit)
}

func (s *Service) CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error) {
	log := logger.Get().WithService("bucket-service")

	if b.ID != "" {
		log.Error("Bucket ID should not be set during creation", "provided_id", b.ID)
		return nil, apperrors.NewValidationError("bucket ID should not be set during creation", nil)
	}

	if b.Key == "" {
		return nil, apperrors.NewValidationError("key is required", nil)
	}

	existing, err := s.Repo.GetBucketByKey(ctx, b.Key)
	if err == nil && existing != nil {
		log.Error("Bucket key already exists", "key", b.Key, "existing_bucket_id", existing.ID)
		return nil, apperrors.NewConflictError("key must be unique", nil)
	}

	if err := s.validateBucket(b); err != nil {
		log.WithError(err).Error("Bucket validation failed")
		return nil, err
	}

	b.ID = generateID()
	b.Status = BucketStatusDraft

	if err := s.Repo.CreateBucket(ctx, b); err != nil {
		log.WithError(err).Error("Failed to create bucket", "bucket_id", b.ID, "key", b.Key)
		return nil, apperrors.NewInternalError("failed to create bucket", err)
	}

	log.Info("Bucket created successfully", "bucket_id", b.ID, "key", b.Key)
	return b, nil
}

func (s *Service) UpdateBucket(ctx context.Context, id string, b *Bucket) error {
	log := logger.Get().WithService("bucket-service")

	existing, err := s.Repo.GetBucketByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("bucket not found", err)
	}

	b.ID = existing.ID
	b.Key = existing.Key
	b.CreatedAt = existing.CreatedAt

	if err := s.validateBucket(b); err != nil {
		log.WithError(err).Error("Bucket validation failed", "bucket_id", id)
		return err
	}

	if err := s.Repo.UpdateBucket(ctx, b); err != nil {
		log.WithError(err).Error("Failed to update bucket", "bucket_id", id)
		return apperrors.NewInternalError("failed to update bucket", err)
	}

	log.Info("Bucket updated successfully", "bucket_id", id)
	return nil
}

func (s *Service) PatchBucket(ctx context.Context, id string, patch map[string]interface{}) error {
	log := logger.Get().WithService("bucket-service")

	if err := s.Repo.PatchBucket(ctx, id, patch); err != nil {
		log.WithError(err).Error("Failed to patch bucket", "bucket_id", id)
		return apperrors.NewInternalError("failed to patch bucket", err)
	}

	log.Info("Bucket patched successfully", "bucket_id", id)
	return nil
}

func (s *Service) DeleteBucket(ctx context.Context, id string) error {
	log := logger.Get().WithService("bucket-service")

	if err := s.Repo.DeleteBucket(ctx, id); err != nil {
		log.WithError(err).Error("Failed to delete bucket", "bucket_id", id)
		return apperrors.NewInternalError("failed to delete bucket", err)
	}

	log.Info("Bucket deleted successfully", "bucket_id", id)
	return nil
}

func (s *Service) AddAssetToBucket(ctx context.Context, bucketID, assetID string) error {
	log := logger.Get().WithService("bucket-service")

	bucket, err := s.GetBucketByID(ctx, bucketID)
	if err != nil {
		return apperrors.WithContext(err, map[string]interface{}{
			"operation": "add_asset_to_bucket",
			"bucket_id": bucketID,
			"asset_id":  assetID,
		})
	}

	for _, existingAssetID := range bucket.AssetIDs {
		if existingAssetID == assetID {
			log.Warn("Asset already exists in bucket", "bucket_id", bucketID, "asset_id", assetID)
			return apperrors.NewConflictError("asset already exists in bucket", nil)
		}
	}

	bucket.AssetIDs = append(bucket.AssetIDs, assetID)

	if err := s.Repo.UpdateBucket(ctx, bucket); err != nil {
		log.WithError(err).Error("Failed to update bucket with new asset", "bucket_id", bucketID, "asset_id", assetID)
		return apperrors.NewInternalError("failed to add asset to bucket", err)
	}

	log.Info("Asset added to bucket successfully", "bucket_id", bucketID, "asset_id", assetID)
	return nil
}

func (s *Service) RemoveAssetFromBucket(ctx context.Context, bucketID, assetID string) error {
	log := logger.Get().WithService("bucket-service")

	bucket, err := s.GetBucketByID(ctx, bucketID)
	if err != nil {
		return apperrors.WithContext(err, map[string]interface{}{
			"operation": "remove_asset_from_bucket",
			"bucket_id": bucketID,
			"asset_id":  assetID,
		})
	}

	found := false
	filtered := make([]string, 0, len(bucket.AssetIDs))
	for _, existingAssetID := range bucket.AssetIDs {
		if existingAssetID == assetID {
			found = true
		} else {
			filtered = append(filtered, existingAssetID)
		}
	}

	if !found {
		log.Warn("Asset not found in bucket", "bucket_id", bucketID, "asset_id", assetID)
		return apperrors.NewNotFoundError("asset not found in bucket", nil)
	}

	bucket.AssetIDs = filtered

	if err := s.Repo.UpdateBucket(ctx, bucket); err != nil {
		log.WithError(err).Error("Failed to update bucket after removing asset", "bucket_id", bucketID, "asset_id", assetID)
		return apperrors.NewInternalError("failed to remove asset from bucket", err)
	}

	log.Info("Asset removed from bucket successfully", "bucket_id", bucketID, "asset_id", assetID)
	return nil
}

func (s *Service) GetBucketsByType(ctx context.Context, bucketType string) ([]Bucket, error) {
	log := logger.Get().WithService("bucket-service")

	if bucketType == "" {
		return nil, apperrors.NewValidationError("bucket type is required", nil)
	}

	buckets, err := s.Repo.GetBucketsByType(ctx, bucketType)
	if err != nil {
		log.WithError(err).Error("Failed to get buckets by type", "bucket_type", bucketType)
		return nil, apperrors.NewInternalError("failed to get buckets by type", err)
	}

	log.Debug("Retrieved buckets by type", "bucket_type", bucketType, "count", len(buckets))
	return buckets, nil
}

func (s *Service) GetBucketsByAsset(ctx context.Context, assetID string) ([]Bucket, error) {
	log := logger.Get().WithService("bucket-service")

	if assetID == "" {
		return nil, apperrors.NewValidationError("asset ID is required", nil)
	}

	buckets, err := s.Repo.GetBucketsByAsset(ctx, assetID)
	if err != nil {
		log.WithError(err).Error("Failed to get buckets by asset", "asset_id", assetID)
		return nil, apperrors.NewInternalError("failed to get buckets by asset", err)
	}

	log.Debug("Retrieved buckets by asset", "asset_id", assetID, "count", len(buckets))
	return buckets, nil
}

func (s *Service) GetAssetsInBucket(ctx context.Context, bucketID string) ([]string, error) {
	log := logger.Get().WithService("bucket-service")

	if bucketID == "" {
		return nil, apperrors.NewValidationError("bucket ID is required", nil)
	}

	assetIDs, err := s.Repo.GetAssetsInBucket(ctx, bucketID)
	if err != nil {
		log.WithError(err).Error("Failed to get assets in bucket", "bucket_id", bucketID)
		return nil, apperrors.NewInternalError("failed to get assets in bucket", err)
	}

	log.Debug("Retrieved assets in bucket", "bucket_id", bucketID, "count", len(assetIDs))
	return assetIDs, nil
}

func (s *Service) validateBucket(b *Bucket) error {
	if b.Name == "" {
		return apperrors.NewValidationError("name is required", nil)
	}
	if b.Type == "" {
		return apperrors.NewValidationError("type is required", nil)
	}
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func NewService(repo BucketRepository) *Service {
	return &Service{Repo: repo}
}
