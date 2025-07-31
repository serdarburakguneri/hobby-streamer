package neo4j

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/infrastructure/neo4j/asset"
)

type AssetRepositoryAdapter struct {
	repo *asset.Repository
}

func NewAssetRepositoryAdapter(repo *asset.Repository) *AssetRepositoryAdapter {
	return &AssetRepositoryAdapter{repo: repo}
}

func (a *AssetRepositoryAdapter) Save(ctx context.Context, asset *entity.Asset) error {
	return a.repo.Save(ctx, asset)
}

func (a *AssetRepositoryAdapter) FindByID(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
	return a.repo.FindByID(ctx, id)
}

func (a *AssetRepositoryAdapter) FindBySlug(ctx context.Context, slug valueobjects.Slug) (*entity.Asset, error) {
	return a.repo.FindBySlug(ctx, slug)
}

func (a *AssetRepositoryAdapter) Update(ctx context.Context, asset *entity.Asset) error {
	return a.repo.Update(ctx, asset)
}

func (a *AssetRepositoryAdapter) Delete(ctx context.Context, id valueobjects.AssetID) error {
	return a.repo.Delete(ctx, id)
}

func (a *AssetRepositoryAdapter) List(ctx context.Context, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.List(ctx, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *AssetRepositoryAdapter) Search(ctx context.Context, query string, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.Search(ctx, query, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *AssetRepositoryAdapter) FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.FindByOwnerID(ctx, ownerID, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *AssetRepositoryAdapter) FindByParentID(ctx context.Context, parentID valueobjects.AssetID, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.FindByParentID(ctx, parentID, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *AssetRepositoryAdapter) FindByType(ctx context.Context, assetType valueobjects.AssetType, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.FindByType(ctx, assetType, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *AssetRepositoryAdapter) FindByGenre(ctx context.Context, genre valueobjects.Genre, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.FindByGenre(ctx, genre, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *AssetRepositoryAdapter) FindByTag(ctx context.Context, tag valueobjects.Tag, limit *int, offset *int) ([]*entity.Asset, error) {
	l, params := toPageParams(limit, offset)
	page, err := a.repo.FindByTag(ctx, tag, l, params)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}
