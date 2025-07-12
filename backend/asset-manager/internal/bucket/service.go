package bucket

import (
	"context"
	"errors"
	"time"
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
	return s.Repo.GetBucketByID(ctx, id)
}

func (s *Service) GetBucketByKey(ctx context.Context, key string) (*Bucket, error) {
	return s.Repo.GetBucketByKey(ctx, key)
}

func (s *Service) ListBuckets(ctx context.Context, limit int) (*BucketPage, error) {
	return s.Repo.ListBuckets(ctx, limit)
}

func (s *Service) CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error) {
	if b.ID != "" {
		return nil, errors.New(ErrIDShouldNotBeSet)
	}

	if b.Key == "" {
		return nil, errors.New("key is required")
	}
	if existing, _ := s.Repo.GetBucketByKey(ctx, b.Key); existing != nil {
		return nil, errors.New("key must be unique")
	}

	if err := s.validateBucket(b); err != nil {
		return nil, err
	}

	b.ID = generateID()
	b.Status = BucketStatusDraft

	if err := s.Repo.CreateBucket(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) UpdateBucket(ctx context.Context, id string, b *Bucket) error {
	existing, err := s.Repo.GetBucketByID(ctx, id)
	if err != nil {
		return err
	}

	b.ID = existing.ID
	b.Key = existing.Key
	b.CreatedAt = existing.CreatedAt

	if err := s.validateBucket(b); err != nil {
		return err
	}

	return s.Repo.UpdateBucket(ctx, b)
}

func (s *Service) PatchBucket(ctx context.Context, id string, patch map[string]interface{}) error {
	return s.Repo.PatchBucket(ctx, id, patch)
}

func (s *Service) DeleteBucket(ctx context.Context, id string) error {
	return s.Repo.DeleteBucket(ctx, id)
}

func (s *Service) AddAssetToBucket(ctx context.Context, bucketID, assetID string) error {
	bucket, err := s.GetBucketByID(ctx, bucketID)
	if err != nil {
		return err
	}
	for _, existingAssetID := range bucket.AssetIDs {
		if existingAssetID == assetID {
			return errors.New(ErrAssetExists)
		}
	}
	bucket.AssetIDs = append(bucket.AssetIDs, assetID)
	return s.Repo.UpdateBucket(ctx, bucket)
}

func (s *Service) RemoveAssetFromBucket(ctx context.Context, bucketID, assetID string) error {
	bucket, err := s.GetBucketByID(ctx, bucketID)
	if err != nil {
		return err
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
		return errors.New(ErrAssetNotFound)
	}
	bucket.AssetIDs = filtered
	return s.Repo.UpdateBucket(ctx, bucket)
}

func (s *Service) GetBucketsByType(ctx context.Context, bucketType string) ([]Bucket, error) {

	// For now, return empty slice
	return []Bucket{}, nil
}

func (s *Service) GetBucketsByAsset(ctx context.Context, assetID string) ([]Bucket, error) {

	// For now, return empty slice
	return []Bucket{}, nil
}

func (s *Service) GetAssetsInBucket(ctx context.Context, bucketID string) ([]string, error) {

	// For now, return empty slice
	return []string{}, nil
}

func (s *Service) validateBucket(b *Bucket) error {
	if b.Name == "" {
		return errors.New("name is required")
	}
	if b.Type == "" {
		return errors.New("type is required")
	}
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func NewService(repo BucketRepository) *Service {
	return &Service{Repo: repo}
}
