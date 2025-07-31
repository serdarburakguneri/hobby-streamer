package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type DomainService interface {
	ValidateAssetHierarchy(ctx context.Context, a *entity.Asset, parentID *valueobjects.AssetID) error
}

type domainService struct {
	repo Repository
}

func NewDomainService(repo Repository) DomainService {
	return &domainService{repo: repo}
}

func (s *domainService) ValidateAssetHierarchy(ctx context.Context, a *entity.Asset, parentID *valueobjects.AssetID) error {
	if parentID == nil {
		return nil
	}
	if a.ID().Value() == parentID.Value() {
		return errors.NewValidationError("asset cannot be its own parent", nil)
	}
	parent, err := s.repo.FindByID(ctx, *parentID)
	if err != nil {
		return errors.NewInternalError("parent asset not found", err)
	}
	if parent == nil {
		return errors.NewNotFoundError("parent asset not found", nil)
	}
	if !parent.IsReadyForPublishing() {
		return errors.NewValidationError("parent asset not published", nil)
	}
	return nil
}
