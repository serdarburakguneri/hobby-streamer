package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AssetRepository interface {
	GetAssetByID(ctx context.Context, id string) (*Asset, error)
	GetAssetBySlug(ctx context.Context, slug string) (*Asset, error)
	GetAssetsByIDs(ctx context.Context, ids []string) ([]Asset, error)
	ListAssets(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	SearchAssets(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	SaveAsset(ctx context.Context, asset *Asset) error
	PatchAsset(ctx context.Context, id string, patch map[string]interface{}) error
	DeleteAsset(ctx context.Context, id string) error
	GetParent(ctx context.Context, childID string) (*Asset, error)
	GetChildren(ctx context.Context, parentID string) ([]Asset, error)
	GetAssetsByTypeAndGenre(ctx context.Context, assetType, genre string) ([]Asset, error)
}

type Repository struct {
	driver neo4j.Driver
	logger *logger.Logger
}

func NewRepository(driver neo4j.Driver) *Repository {
	return &Repository{
		driver: driver,
		logger: logger.WithService("asset-repository"),
	}
}

func (r *Repository) SaveAsset(ctx context.Context, a *Asset) error {
	log := r.logger.WithContext(ctx)

	now := time.Now().UTC()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MERGE (a:Asset {id: $id})
		ON CREATE SET
			a.slug = $slug,
					a.title = $title,
		a.description = $description,
		a.type = $type,
		a.genre = $genre,
		a.genres = $genres,
		a.tags = $tags,
		a.ownerId = $ownerId,
		a.publishRulePublishAt = $publishRulePublishAt,
		a.publishRuleUnpublishAt = $publishRuleUnpublishAt,
		a.publishRuleRegions = $publishRuleRegions,
		a.publishRuleAgeRating = $publishRuleAgeRating,
		a.videos = $videos,
		a.createdAt = $createdAt,
		a.updatedAt = $updatedAt
	ON MATCH SET
		a.title = $title,
		a.description = $description,
		a.type = $type,
		a.genre = $genre,
		a.genres = $genres,
		a.tags = $tags,
		a.ownerId = $ownerId,
		a.publishRulePublishAt = $publishRulePublishAt,
		a.publishRuleUnpublishAt = $publishRuleUnpublishAt,
		a.publishRuleRegions = $publishRuleRegions,
		a.publishRuleAgeRating = $publishRuleAgeRating,
		a.videos = $videos,
		a.updatedAt = $updatedAt
		RETURN a
	`

	var title, description, assetType, genre, ownerID string
	if a.Title != nil {
		title = *a.Title
	}
	if a.Description != nil {
		description = *a.Description
	}
	if a.Type != nil {
		assetType = *a.Type
	}
	if a.Genre != nil {
		genre = *a.Genre
	}
	if a.OwnerID != nil {
		ownerID = *a.OwnerID
	}

	var publishRulePublishAt time.Time
	var publishRuleUnpublishAt time.Time
	var publishRuleRegions []string
	var publishRuleAgeRating string

	if a.PublishRule != nil {
		publishRulePublishAt = a.PublishRule.PublishAt
		publishRuleUnpublishAt = a.PublishRule.UnpublishAt
		publishRuleRegions = a.PublishRule.Regions
		publishRuleAgeRating = a.PublishRule.AgeRating
	}

	videosJSON, err := json.Marshal(a.Videos)
	if err != nil {
		log.WithError(err).Error("Failed to marshal videos to JSON", "asset_id", a.ID)
		return fmt.Errorf("failed to marshal videos: %w", err)
	}

	params := map[string]interface{}{
		"id":                     a.ID,
		"slug":                   a.Slug,
		"title":                  title,
		"description":            description,
		"type":                   assetType,
		"genre":                  genre,
		"genres":                 a.Genres,
		"tags":                   a.Tags,
		"ownerId":                ownerID,
		"publishRulePublishAt":   publishRulePublishAt,
		"publishRuleUnpublishAt": publishRuleUnpublishAt,
		"publishRuleRegions":     publishRuleRegions,
		"publishRuleAgeRating":   publishRuleAgeRating,
		"videos":                 string(videosJSON),
		"createdAt":              a.CreatedAt,
		"updatedAt":              a.UpdatedAt,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to save asset to Neo4j", "asset_id", a.ID)
		return fmt.Errorf("failed to save asset: %w", err)
	}

	if !result.Next() {
		log.Error("Asset save query did not return any results", "asset_id", a.ID)
		return fmt.Errorf("asset save query did not return any results")
	}

	// Handle parent relationship if parentId is provided
	if a.ParentID != nil && *a.ParentID != "" {
		parentQuery := `
			MATCH (child:Asset {id: $childId})
			MATCH (parent:Asset {id: $parentId})
			MERGE (child)-[:BELONGS_TO]->(parent)
			RETURN child, parent
		`
		parentParams := map[string]interface{}{
			"childId":  a.ID,
			"parentId": *a.ParentID,
		}

		_, err = session.Run(parentQuery, parentParams)
		if err != nil {
			log.WithError(err).Error("Failed to create parent relationship", "asset_id", a.ID, "parent_id", *a.ParentID)
			return fmt.Errorf("failed to create parent relationship: %w", err)
		}
		log.Debug("Parent relationship created successfully", "asset_id", a.ID, "parent_id", *a.ParentID)
	}

	return nil
}

func (r *Repository) GetAssetByID(ctx context.Context, id string) (*Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {id: $id})
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
	`

	result, err := session.Run(query, map[string]interface{}{"id": id})
	if err != nil {
		log.WithError(err).Error("Failed to get asset from Neo4j", "asset_id", id)
		return nil, fmt.Errorf("get asset failed: %w", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Debug("Asset not found in Neo4j", "asset_id", id)
		return nil, fmt.Errorf("asset not found")
	}

	asset, err := r.recordToAssetWithParent(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to asset", "asset_id", id)
		return nil, fmt.Errorf("convert record to asset failed: %w", err)
	}

	log.Debug("Asset retrieved successfully from Neo4j", "asset_id", id, "title", asset.Title)
	return asset, nil
}

func (r *Repository) GetAssetBySlug(ctx context.Context, slug string) (*Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {slug: $slug})
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
	`

	result, err := session.Run(query, map[string]interface{}{"slug": slug})
	if err != nil {
		log.WithError(err).Error("Failed to get asset by slug from Neo4j", "slug", slug)
		return nil, fmt.Errorf("get asset by slug failed: %w", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Debug("Asset not found in Neo4j", "slug", slug)
		return nil, fmt.Errorf("asset not found")
	}

	asset, err := r.recordToAssetWithParent(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to asset", "slug", slug)
		return nil, fmt.Errorf("convert record to asset failed: %w", err)
	}

	log.Debug("Asset retrieved successfully from Neo4j", "slug", slug, "title", asset.Title)
	return asset, nil
}

func (r *Repository) ListAssets(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset)
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to list assets from Neo4j")
		return nil, fmt.Errorf("list assets failed: %w", err)
	}

	var assets []Asset
	for result.Next() {
		asset, err := r.recordToAssetWithParent(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to asset, skipping")
			continue
		}
		assets = append(assets, *asset)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j results")
		return nil, fmt.Errorf("iterate results failed: %w", err)
	}

	return &AssetPage{
		Items: assets,
	}, nil
}

func (r *Repository) SearchAssets(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	cypherQuery := `
		MATCH (a:Asset)
		WHERE a.title CONTAINS $query
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"query": query,
		"limit": limit,
	}

	result, err := session.Run(cypherQuery, params)
	if err != nil {
		log.WithError(err).Error("Failed to search assets from Neo4j", "query", query)
		return nil, fmt.Errorf("search assets failed: %w", err)
	}

	var assets []Asset
	for result.Next() {
		asset, err := r.recordToAssetWithParent(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to asset, skipping")
			continue
		}
		assets = append(assets, *asset)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j search results")
		return nil, fmt.Errorf("iterate search results failed: %w", err)
	}

	log.Debug("Search completed successfully", "query", query, "results_count", len(assets))
	return &AssetPage{
		Items: assets,
	}, nil
}

func (r *Repository) PatchAsset(ctx context.Context, id string, patch map[string]interface{}) error {
	log := r.logger.WithContext(ctx)

	if len(patch) == 0 {
		log.Debug("No patch fields provided", "asset_id", id)
		return nil
	}

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// Build dynamic SET clause
	setClause := "SET a.updatedAt = $updatedAt"
	params := map[string]interface{}{
		"id":        id,
		"updatedAt": time.Now().UTC(),
	}

	// Track if parentId is being updated
	var parentIdUpdated bool
	var newParentId string

	for key, value := range patch {
		paramKey := fmt.Sprintf("param_%s", key)
		setClause += fmt.Sprintf(", a.%s = $%s", key, paramKey)
		params[paramKey] = value

		// Track parentId updates
		if key == "parentId" {
			parentIdUpdated = true
			if value != nil {
				newParentId = value.(string)
			}
		}
	}

	query := fmt.Sprintf(`
		MATCH (a:Asset {id: $id})
		%s
		RETURN a
	`, setClause)

	_, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to patch asset in Neo4j", "asset_id", id)
		return fmt.Errorf("failed to patch asset: %w", err)
	}

	if parentIdUpdated {

		removeParentQuery := `
			MATCH (child:Asset {id: $childId})-[r:BELONGS_TO]->(parent:Asset)
			DELETE r
		`
		_, err = session.Run(removeParentQuery, map[string]interface{}{"childId": id})
		if err != nil {
			log.WithError(err).Warn("Failed to remove existing parent relationship", "asset_id", id)
		}

		// Then, create new parent relationship if parentId is not empty
		if newParentId != "" {
			parentQuery := `
				MATCH (child:Asset {id: $childId})
				MATCH (parent:Asset {id: $parentId})
				MERGE (child)-[:BELONGS_TO]->(parent)
				RETURN child, parent
			`
			parentParams := map[string]interface{}{
				"childId":  id,
				"parentId": newParentId,
			}

			_, err = session.Run(parentQuery, parentParams)
			if err != nil {
				log.WithError(err).Error("Failed to create parent relationship", "asset_id", id, "parent_id", newParentId)
				return fmt.Errorf("failed to create parent relationship: %w", err)
			}
			log.Debug("Parent relationship updated successfully", "asset_id", id, "parent_id", newParentId)
		} else {
			log.Debug("Parent relationship removed successfully", "asset_id", id)
		}
	}

	log.Debug("Asset patched successfully in Neo4j", "asset_id", id, "fields_updated", len(patch))
	return nil
}

func (r *Repository) recordToAssetWithParent(record *neo4j.Record) (*Asset, error) {
	asset, err := r.recordToAsset(record)
	if err != nil {
		return nil, err
	}

	// Check if parent exists in the record
	if parentNode, ok := record.Get("parent"); ok && parentNode != nil {
		if parentNeo4jNode, ok := parentNode.(neo4j.Node); ok {
			parentProps := parentNeo4jNode.Props

			parentAsset := &Asset{}
			if id, ok := parentProps["id"].(string); ok {
				parentAsset.ID = id
			}
			if title, ok := parentProps["title"].(string); ok {
				parentAsset.Title = &title
			}
			if assetType, ok := parentProps["type"].(string); ok {
				parentAsset.Type = &assetType
			}

			asset.Parent = parentAsset
			if parentAsset.ID != "" {
				asset.ParentID = &parentAsset.ID
			}
		}
	}

	return asset, nil
}

func (r *Repository) recordToAsset(record *neo4j.Record) (*Asset, error) {
	node, ok := record.Get("a")
	if !ok {
		return nil, fmt.Errorf("no 'a' field in record")
	}

	neo4jNode, ok := node.(neo4j.Node)
	if !ok {
		return nil, fmt.Errorf("field 'a' is not a node")
	}

	props := neo4jNode.Props

	asset := &Asset{}

	// Required fields with type assertions
	if id, ok := props["id"].(string); ok {
		asset.ID = id
	}
	if slug, ok := props["slug"].(string); ok {
		asset.Slug = slug
	}
	if createdAt, ok := props["createdAt"].(time.Time); ok {
		asset.CreatedAt = createdAt
	} else if createdAtStr, ok := props["createdAt"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			asset.CreatedAt = parsed
		} else {
			asset.CreatedAt = time.Now()
		}
	} else {
		asset.CreatedAt = time.Now()
	}

	if updatedAt, ok := props["updatedAt"].(time.Time); ok {
		asset.UpdatedAt = updatedAt
	} else if updatedAtStr, ok := props["updatedAt"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			asset.UpdatedAt = parsed
		} else {
			asset.UpdatedAt = time.Now()
		}
	} else {
		asset.UpdatedAt = time.Now()
	}

	if title, ok := props["title"].(string); ok {
		asset.Title = &title
	}
	if description, ok := props["description"].(string); ok {
		asset.Description = &description
	}

	if assetType, ok := props["type"].(string); ok {
		asset.Type = &assetType
	}
	if genre, ok := props["genre"].(string); ok {
		asset.Genre = &genre
	}
	if ownerID, ok := props["ownerId"].(string); ok {
		asset.OwnerID = &ownerID
	}

	if genres, ok := props["genres"].([]interface{}); ok {
		for _, g := range genres {
			if genreStr, ok := g.(string); ok {
				asset.Genres = append(asset.Genres, genreStr)
			}
		}
	}
	if tags, ok := props["tags"].([]interface{}); ok {
		for _, t := range tags {
			if tagStr, ok := t.(string); ok {
				asset.Tags = append(asset.Tags, tagStr)
			}
		}
	}

	if metadata, ok := props["attributes"].(map[string]interface{}); ok {
		asset.Metadata = make(map[string]interface{})
		for k, v := range metadata {
			asset.Metadata[k] = v
		}
	}

	if props["publishRulePublishAt"] != nil {
		publishRule := &PublishRule{}

		if publishAt, ok := props["publishRulePublishAt"].(time.Time); ok {
			publishRule.PublishAt = publishAt
		} else if publishAtStr, ok := props["publishRulePublishAt"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, publishAtStr); err == nil {
				publishRule.PublishAt = parsed
			}
		}

		if unpublishAt, ok := props["publishRuleUnpublishAt"].(time.Time); ok {
			publishRule.UnpublishAt = unpublishAt
		} else if unpublishAtStr, ok := props["publishRuleUnpublishAt"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, unpublishAtStr); err == nil {
				publishRule.UnpublishAt = parsed
			}
		}

		if regions, ok := props["publishRuleRegions"].([]interface{}); ok {
			for _, r := range regions {
				if regionStr, ok := r.(string); ok {
					publishRule.Regions = append(publishRule.Regions, regionStr)
				}
			}
		}

		if ageRating, ok := props["publishRuleAgeRating"].(string); ok {
			publishRule.AgeRating = ageRating
		}

		asset.PublishRule = publishRule
	}

	if videosJSON, ok := props["videos"].(string); ok && videosJSON != "" {
		var videos []Video
		if err := json.Unmarshal([]byte(videosJSON), &videos); err != nil {
			r.logger.WithError(err).Warn("Failed to unmarshal videos JSON, setting empty videos array", "asset_id", asset.ID)
			asset.Videos = []Video{}
		} else {
			asset.Videos = videos
		}
	} else {
		asset.Videos = []Video{}
	}

	return asset, nil
}

func (r *Repository) GetChildren(ctx context.Context, parentID string) ([]Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (child:Asset)-[:BELONGS_TO]->(parent:Asset {id: $parentID})
		RETURN child as a
		ORDER BY child.createdAt
	`

	result, err := session.Run(query, map[string]interface{}{"parentID": parentID})
	if err != nil {
		log.WithError(err).Error("Failed to get children from Neo4j", "parent_id", parentID)
		return nil, fmt.Errorf("get children failed: %w", err)
	}

	var assets []Asset
	for result.Next() {
		asset, err := r.recordToAsset(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to child asset, skipping")
			continue
		}
		assets = append(assets, *asset)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j children results")
		return nil, fmt.Errorf("iterate children results failed: %w", err)
	}

	log.Debug("Children retrieved successfully from Neo4j", "parent_id", parentID, "count", len(assets))
	return assets, nil
}

func (r *Repository) GetParent(ctx context.Context, childID string) (*Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (child:Asset {id: $childID})-[:BELONGS_TO]->(parent:Asset)
		RETURN parent
	`

	result, err := session.Run(query, map[string]interface{}{"childID": childID})
	if err != nil {
		log.WithError(err).Error("Failed to get parent from Neo4j", "child_id", childID)
		return nil, fmt.Errorf("get parent failed: %w", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Debug("Parent not found in Neo4j", "child_id", childID)
		return nil, nil
	}

	parent, err := r.recordToAsset(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to parent asset", "child_id", childID)
		return nil, fmt.Errorf("convert record to parent asset failed: %w", err)
	}

	log.Debug("Parent retrieved successfully from Neo4j", "child_id", childID, "parent_id", parent.ID)
	return parent, nil
}

func (r *Repository) GetAssetsByTypeAndGenre(ctx context.Context, assetType, genre string) ([]Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {type: $type})
		WHERE $genre IN a.genres OR a.genre = $genre
		RETURN a
		ORDER BY a.createdAt DESC
	`

	params := map[string]interface{}{
		"type":  assetType,
		"genre": genre,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to get assets by type and genre from Neo4j", "type", assetType, "genre", genre)
		return nil, fmt.Errorf("get assets by type and genre failed: %w", err)
	}

	var assets []Asset
	for result.Next() {
		asset, err := r.recordToAsset(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to asset, skipping")
			continue
		}
		assets = append(assets, *asset)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j type/genre results")
		return nil, fmt.Errorf("iterate type/genre results failed: %w", err)
	}

	log.Debug("Assets by type and genre retrieved successfully from Neo4j", "type", assetType, "genre", genre, "count", len(assets))
	return assets, nil
}

func (r *Repository) DeleteAsset(ctx context.Context, id string) error {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {id: $id})
		DETACH DELETE a
	`

	result, err := session.Run(query, map[string]interface{}{"id": id})
	if err != nil {
		log.WithError(err).Error("Failed to delete asset from Neo4j", "asset_id", id)
		return fmt.Errorf("delete asset failed: %w", err)
	}

	summary, err := result.Consume()
	if err != nil {
		log.WithError(err).Error("Failed to consume delete result", "asset_id", id)
		return fmt.Errorf("consume delete result failed: %w", err)
	}

	if summary.Counters().NodesDeleted() == 0 {
		log.Debug("Asset not found for deletion", "asset_id", id)
		return fmt.Errorf("asset not found")
	}

	log.Debug("Asset deleted successfully from Neo4j", "asset_id", id)
	return nil
}

func (r *Repository) GetAssetsByIDs(ctx context.Context, ids []string) ([]Asset, error) {
	log := r.logger.WithContext(ctx)

	if len(ids) == 0 {
		return []Asset{}, nil
	}

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset)
		WHERE a.id IN $ids
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
		ORDER BY a.id
	`

	result, err := session.Run(query, map[string]interface{}{"ids": ids})
	if err != nil {
		log.WithError(err).Error("Failed to get assets by IDs from Neo4j", "asset_ids", ids)
		return nil, fmt.Errorf("get assets by IDs failed: %w", err)
	}

	var assets []Asset
	for result.Next() {
		asset, err := r.recordToAssetWithParent(result.Record())
		if err != nil {
			log.WithError(err).Error("Failed to parse asset record", "asset_ids", ids)
			return nil, fmt.Errorf("failed to parse asset record: %w", err)
		}
		assets = append(assets, *asset)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating over assets result", "asset_ids", ids)
		return nil, fmt.Errorf("error iterating over assets result: %w", err)
	}

	log.Debug("Retrieved assets by IDs", "count", len(assets), "requested_ids", len(ids))
	return assets, nil
}
