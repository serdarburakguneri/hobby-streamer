package asset

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Define an AssetRepository interface for service injection
type AssetRepository interface {
	SaveAsset(ctx context.Context, a *Asset) error
	GetAssetByID(ctx context.Context, id int) (*Asset, error)
	ListAssets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error)
	PatchAsset(ctx context.Context, id int, patch map[string]interface{}) error
}

type Service struct {
	Repo AssetRepository
}

// Ensure Service implements the AssetService interface for handler injection
var _ AssetService = (*Service)(nil)

func generateID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1_000_000_000)
}

func (s *Service) GetAssetByID(ctx context.Context, id int) (*Asset, error) {
	return s.Repo.GetAssetByID(ctx, id)
}

func (s *Service) ListAssets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error) {
	return s.Repo.ListAssets(ctx, limit, lastKey)
}

func (s *Service) CreateAsset(ctx context.Context, a *Asset) (*Asset, error) {
	if a.ID != 0 {
		return nil, errors.New("id should not be set by client")
	}
	a.ID = generateID()
	if err := s.Repo.SaveAsset(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) UpdateAsset(ctx context.Context, id int, a *Asset) error {
	if a.ID != id {
		return errors.New("id mismatch")
	}
	return s.Repo.SaveAsset(ctx, a)
}

func (s *Service) PatchAsset(ctx context.Context, id int, patch map[string]interface{}) error {
	return s.Repo.PatchAsset(ctx, id, patch)
}

func (s *Service) AddImage(ctx context.Context, id int, img *Image) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	for _, existing := range asset.Images {
		if existing.FileName == img.FileName {
			return errors.New("image with same filename already exists")
		}
	}

	asset.Images = append(asset.Images, *img)
	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) DeleteImage(ctx context.Context, id int, filename string) error {
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

func (s *Service) AddVideo(ctx context.Context, id int, label string, video *Video) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	if asset.Videos == nil {
		asset.Videos = make(map[string]Video)
	}
	asset.Videos[label] = *video
	return s.Repo.SaveAsset(ctx, asset)
}

func (s *Service) DeleteVideo(ctx context.Context, id int, label string) error {
	asset, err := s.Repo.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}

	delete(asset.Videos, label)
	return s.Repo.SaveAsset(ctx, asset)
}

func NewService(repo AssetRepository) *Service {
	return &Service{Repo: repo}
}
