package asset

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Repository struct {
	driver    neo4j.Driver
	logger    *logger.Logger
	converter *AssetConverter
}

func NewRepository(driver neo4j.Driver) *Repository {
	logger := logger.WithService("neo4j-asset-repository")
	return &Repository{
		driver:    driver,
		logger:    logger,
		converter: NewAssetConverter(logger),
	}
}

func (r *Repository) Save(ctx context.Context, a *entity.Asset) error {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetSaveQuery()
	params := r.converter.AssetToParams(a)

	_, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to save asset to Neo4j", "asset_id", a.ID())
		return pkgerrors.NewInternalError("database operation failed: unable to save asset", err)
	}

	if a.ParentID() != nil {
		if err := r.createParentRelationship(session, a.ID().Value(), a.ParentID().Value()); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindByIDQuery()
	params := map[string]interface{}{"id": id.Value()}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to find asset by ID", "asset_id", id.Value())
		return nil, pkgerrors.NewInternalError("database operation failed: unable to retrieve asset", err)
	}

	if !result.Next() {
		return nil, pkgerrors.NewNotFoundError("asset not found with the specified ID", nil)
	}

	return r.converter.RecordToAsset(result.Record())
}

func (r *Repository) FindBySlug(ctx context.Context, slug valueobjects.Slug) (*entity.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindBySlugQuery()

	result, err := session.Run(query, map[string]interface{}{"slug": slug.Value()})
	if err != nil {
		log.WithError(err).Error("Failed to get asset by slug from Neo4j", "slug", slug.Value())
		return nil, pkgerrors.NewInternalError("get asset by slug failed", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Warn("Asset not found in Neo4j", "slug", slug.Value())
		return nil, pkgerrors.NewNotFoundError("asset not found", err)
	}

	return r.converter.RecordToAsset(record)
}

func (r *Repository) Update(ctx context.Context, a *entity.Asset) error {
	err := r.Save(ctx, a)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id valueobjects.AssetID) error {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetDeleteQuery()

	_, err := session.Run(query, map[string]interface{}{"id": id.Value()})
	if err != nil {
		log.WithError(err).Error("Failed to delete asset from Neo4j", "asset_id", id.Value())
		return pkgerrors.NewInternalError("failed to delete asset", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetListQuery()

	params := map[string]interface{}{
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to list assets from Neo4j")
		return nil, pkgerrors.NewInternalError("list assets failed", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) Search(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	queryStr := buildAssetSearchQuery()
	params := map[string]interface{}{
		"query":   query,
		"limit":   limit,
		"lastKey": lastKey,
	}

	result, err := session.Run(queryStr, params)
	if err != nil {
		log.WithError(err).Error("Failed to search assets in Neo4j", "query", query)
		return nil, pkgerrors.NewInternalError("database operation failed: unable to search assets", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert Neo4j record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindByOwnerIDQuery()
	params := map[string]interface{}{
		"ownerId": ownerID.Value(),
		"limit":   limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to find assets by owner ID", "owner_id", ownerID.Value())
		return nil, pkgerrors.NewInternalError("database operation failed: unable to retrieve assets by owner", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) FindByParentID(ctx context.Context, parentID valueobjects.AssetID, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindByParentIDQuery()
	params := map[string]interface{}{
		"parentId": parentID.Value(),
		"limit":    limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to find assets by parent ID", "parent_id", parentID.Value())
		return nil, pkgerrors.NewInternalError("database operation failed: unable to retrieve assets by parent", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) FindByType(ctx context.Context, assetType valueobjects.AssetType, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindByTypeQuery()
	params := map[string]interface{}{
		"type":  assetType.Value(),
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to find assets by type", "type", assetType.Value())
		return nil, pkgerrors.NewInternalError("database operation failed: unable to retrieve assets by type", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) FindByGenre(ctx context.Context, genre valueobjects.Genre, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindByGenreQuery()
	params := map[string]interface{}{
		"genre": genre.Value(),
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to find assets by genre", "genre", genre.Value())
		return nil, pkgerrors.NewInternalError("database operation failed: unable to retrieve assets by genre", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) FindByTag(ctx context.Context, tag valueobjects.Tag, limit int, lastKey map[string]interface{}) (*entity.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := buildAssetFindByTagQuery()
	params := map[string]interface{}{
		"tag":   tag.Value(),
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to find assets by tag", "tag", tag.Value())
		return nil, pkgerrors.NewInternalError("database operation failed: unable to retrieve assets by tag", err)
	}

	var assets []*entity.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.converter.RecordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &entity.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		LastKey: lastKey,
	}, nil
}

func (r *Repository) createParentRelationship(session neo4j.Session, childID, parentID string) error {
	query := buildParentRelationshipQuery()
	params := map[string]interface{}{
		"childID":  childID,
		"parentID": parentID,
	}

	_, err := session.Run(query, params)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create parent relationship", "child_id", childID, "parent_id", parentID)
		return pkgerrors.NewInternalError("database operation failed: unable to create parent-child relationship", err)
	}

	return nil
}
