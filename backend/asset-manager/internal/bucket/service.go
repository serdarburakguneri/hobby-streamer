package bucket

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BucketService interface {
	GetBucketByID(ctx context.Context, id int) (*Bucket, error)
	ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*BucketPage, error)
	CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error)
	UpdateBucket(ctx context.Context, id int, b *Bucket) error
	PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error
	AddAssetToBucket(ctx context.Context, id int, assetID int) error
	RemoveAssetFromBucket(ctx context.Context, id int, assetID int) error
}

type Service struct {
	Repo BucketRepository
}

var _ BucketService = (*Service)(nil)

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(1_000_000_000))
}

func (s *Service) GetBucketByID(ctx context.Context, id int) (*Bucket, error) {
	return s.Repo.GetBucketByID(ctx, id)
}

func (s *Service) ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*BucketPage, error) {
	return s.Repo.ListBuckets(ctx, limit, lastKey)
}

func (s *Service) CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error) {
	if b.ID != "" {
		return nil, errors.New(ErrIDShouldNotBeSet)
	}
	b.ID = generateID()
	if err := s.Repo.SaveBucket(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) UpdateBucket(ctx context.Context, id int, b *Bucket) error {
	idStr := strconv.Itoa(id)
	if b.ID != idStr {
		return errors.New(ErrIDMismatch)
	}
	return s.Repo.SaveBucket(ctx, b)
}

func (s *Service) PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error {
	return s.Repo.PatchBucket(ctx, id, patch)
}

func (s *Service) AddAssetToBucket(ctx context.Context, bucketID int, assetID int) error {
	bucket, err := s.GetBucketByID(ctx, bucketID)
	if err != nil {
		return err
	}

	assetIDStr := strconv.Itoa(assetID)
	for _, existingAssetID := range bucket.AssetIDs {
		if existingAssetID == assetIDStr {
			return errors.New(ErrAssetExists)
		}
	}

	bucket.AssetIDs = append(bucket.AssetIDs, assetIDStr)
	return s.Repo.SaveBucket(ctx, bucket)
}

func (s *Service) RemoveAssetFromBucket(ctx context.Context, bucketID int, assetID int) error {
	bucket, err := s.GetBucketByID(ctx, bucketID)
	if err != nil {
		return err
	}

	assetIDStr := strconv.Itoa(assetID)
	found := false
	filtered := make([]string, 0, len(bucket.AssetIDs))
	for _, existingAssetID := range bucket.AssetIDs {
		if existingAssetID == assetIDStr {
			found = true
		} else {
			filtered = append(filtered, existingAssetID)
		}
	}

	if !found {
		return errors.New(ErrAssetNotFound)
	}

	bucket.AssetIDs = filtered
	return s.Repo.SaveBucket(ctx, bucket)
}

func NewService(repo BucketRepository) *Service {
	return &Service{Repo: repo}
}
