package bucket

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Define a BucketRepository interface for service injection
type BucketRepository interface {
	SaveBucket(ctx context.Context, b *Bucket) error
	GetBucketByID(ctx context.Context, id int) (*Bucket, error)
	ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*BucketPage, error)
	PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error
}

type Service struct {
	Repo BucketRepository
}

func NewService(repo BucketRepository) *Service {
	return &Service{Repo: repo}
}

func generateID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1_000_000_000)
}

func (s *Service) GetBucketByID(ctx context.Context, id int) (*Bucket, error) {
	return s.Repo.GetBucketByID(ctx, id)
}

func (s *Service) ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*BucketPage, error) {
	return s.Repo.ListBuckets(ctx, limit, lastKey)
}

func (s *Service) CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error) {
	if b.ID != "" {
		return nil, errors.New("id should not be set by client")
	}
	b.ID = strconv.Itoa(generateID())
	if err := s.Repo.SaveBucket(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) UpdateBucket(ctx context.Context, id int, b *Bucket) error {
	idStr := strconv.Itoa(id)
	if b.ID != idStr {
		return errors.New("id mismatch")
	}
	return s.Repo.SaveBucket(ctx, b)
}

func (s *Service) PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error {
	return s.Repo.PatchBucket(ctx, id, patch)
}

func (s *Service) AddAssetToBucket(ctx context.Context, bucketID int, assetID int) error {
	bucket, err := s.Repo.GetBucketByID(ctx, bucketID)
	if err != nil {
		return err
	}

	idStr := strconv.Itoa(assetID)
	for _, existing := range bucket.AssetIDs {
		if existing == idStr {
			return errors.New("asset already in bucket")
		}
	}

	bucket.AssetIDs = append(bucket.AssetIDs, idStr)
	return s.Repo.SaveBucket(ctx, bucket)
}

func (s *Service) RemoveAssetFromBucket(ctx context.Context, bucketID int, assetID int) error {
	bucket, err := s.Repo.GetBucketByID(ctx, bucketID)
	if err != nil {
		return err
	}

	idStr := strconv.Itoa(assetID)
	filtered := make([]string, 0, len(bucket.AssetIDs))
	for _, a := range bucket.AssetIDs {
		if a != idStr {
			filtered = append(filtered, a)
		}
	}

	bucket.AssetIDs = filtered
	return s.Repo.SaveBucket(ctx, bucket)
}

// Define a BucketService interface for handler and test injection
//go:generate mockgen -destination=../../test/mocks/bucket_service_mock.go -package=mocks github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/bucket BucketService

type BucketService interface {
	GetBucketByID(ctx context.Context, id int) (*Bucket, error)
	ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*BucketPage, error)
	CreateBucket(ctx context.Context, b *Bucket) (*Bucket, error)
	UpdateBucket(ctx context.Context, id int, b *Bucket) error
	PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error
	AddAssetToBucket(ctx context.Context, id int, assetID int) error
	RemoveAssetFromBucket(ctx context.Context, id int, assetID int) error
}

// Ensure Service implements BucketService
var _ BucketService = (*Service)(nil)
