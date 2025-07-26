package neo4j

import (
	"context"
	"encoding/json"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AssetRepository struct {
	driver neo4j.Driver
	logger *logger.Logger
}

func NewAssetRepository(driver neo4j.Driver) *AssetRepository {
	return &AssetRepository{
		driver: driver,
		logger: logger.WithService("neo4j-asset-repository"),
	}
}

func (r *AssetRepository) Save(ctx context.Context, a *asset.Asset) error {
	log := r.logger.WithContext(ctx)

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
			a.publishRule = $publishRule,
			a.videos = $videos,
			a.images = $images,
			a.credits = $credits,
			a.metadata = $metadata,
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
			a.publishRule = $publishRule,
			a.videos = $videos,
			a.images = $images,
			a.credits = $credits,
			a.metadata = $metadata,
			a.updatedAt = $updatedAt
		RETURN a
	`

	params := r.assetToParams(a)

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to save asset to Neo4j", "asset_id", a.ID())
		return pkgerrors.NewInternalError("failed to save asset", err)
	}

	if !result.Next() {
		log.Error("Asset save query did not return any results", "asset_id", a.ID())
		return pkgerrors.NewInternalError("asset save query did not return any results", nil)
	}

	if a.ParentID() != nil {
		if err := r.createParentRelationship(session, a.ID().Value(), a.ParentID().Value()); err != nil {
			return err
		}
	}

	return nil
}

func (r *AssetRepository) FindByID(ctx context.Context, id asset.AssetID) (*asset.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {id: $id})
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
	`

	result, err := session.Run(query, map[string]interface{}{"id": id.Value()})
	if err != nil {
		log.WithError(err).Error("Failed to get asset from Neo4j", "asset_id", id.Value())
		return nil, pkgerrors.NewInternalError("get asset failed", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Warn("Asset not found in Neo4j", "asset_id", id.Value())
		return nil, pkgerrors.NewNotFoundError("asset not found", err)
	}

	asset, err := r.recordToAsset(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to asset", "asset_id", id.Value())
		return nil, pkgerrors.NewInternalError("convert record to asset failed", err)
	}

	return asset, nil
}

func (r *AssetRepository) FindBySlug(ctx context.Context, slug asset.Slug) (*asset.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {slug: $slug})
		OPTIONAL MATCH (a)-[:BELONGS_TO]->(parent:Asset)
		RETURN a, parent
	`

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

	asset, err := r.recordToAsset(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to asset", "slug", slug.Value())
		return nil, pkgerrors.NewInternalError("convert record to asset failed", err)
	}

	return asset, nil
}

func (r *AssetRepository) Update(ctx context.Context, a *asset.Asset) error {
	return r.Save(ctx, a)
}

func (r *AssetRepository) Delete(ctx context.Context, id asset.AssetID) error {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset {id: $id})
		OPTIONAL MATCH (a)-[r]-()
		DELETE r, a
	`

	_, err := session.Run(query, map[string]interface{}{"id": id.Value()})
	if err != nil {
		log.WithError(err).Error("Failed to delete asset from Neo4j", "asset_id", id.Value())
		return pkgerrors.NewInternalError("failed to delete asset", err)
	}

	return nil
}

func (r *AssetRepository) List(ctx context.Context, limit int, lastKey map[string]interface{}) (*asset.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset)
		RETURN a
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to list assets from Neo4j")
		return nil, pkgerrors.NewInternalError("list assets failed", err)
	}

	var assets []*asset.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.recordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &asset.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		Total:   len(assets),
	}, nil
}

func (r *AssetRepository) Search(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*asset.AssetPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	cypherQuery := `
		MATCH (a:Asset)
		WHERE a.title CONTAINS $query OR a.description CONTAINS $query OR a.slug CONTAINS $query
		RETURN a
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"query": query,
		"limit": limit,
	}

	result, err := session.Run(cypherQuery, params)
	if err != nil {
		log.WithError(err).Error("Failed to search assets from Neo4j")
		return nil, pkgerrors.NewInternalError("search assets failed", err)
	}

	var assets []*asset.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.recordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return &asset.AssetPage{
		Items:   assets,
		HasMore: len(assets) == limit,
		Total:   len(assets),
	}, nil
}

func (r *AssetRepository) FindByIDs(ctx context.Context, ids []asset.AssetID) ([]*asset.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.Value()
	}

	query := `
		MATCH (a:Asset)
		WHERE a.id IN $ids
		RETURN a
	`

	params := map[string]interface{}{
		"ids": idStrings,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to get assets by IDs from Neo4j")
		return nil, pkgerrors.NewInternalError("get assets by IDs failed", err)
	}

	var assets []*asset.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.recordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *AssetRepository) FindParent(ctx context.Context, childID asset.AssetID) (*asset.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (child:Asset {id: $childId})-[:BELONGS_TO]->(parent:Asset)
		RETURN parent
	`

	result, err := session.Run(query, map[string]interface{}{"childId": childID.Value()})
	if err != nil {
		log.WithError(err).Error("Failed to get parent asset from Neo4j", "child_id", childID.Value())
		return nil, pkgerrors.NewInternalError("get parent asset failed", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Debug("Parent asset not found in Neo4j", "child_id", childID.Value())
		return nil, pkgerrors.NewNotFoundError("parent asset not found", err)
	}

	asset, err := r.recordToAsset(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to asset", "child_id", childID.Value())
		return nil, pkgerrors.NewInternalError("convert record to asset failed", err)
	}

	return asset, nil
}

func (r *AssetRepository) FindChildren(ctx context.Context, parentID asset.AssetID) ([]*asset.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (parent:Asset {id: $parentId})<-[:BELONGS_TO]-(child:Asset)
		RETURN child
	`

	result, err := session.Run(query, map[string]interface{}{"parentId": parentID.Value()})
	if err != nil {
		log.WithError(err).Error("Failed to get children assets from Neo4j", "parent_id", parentID.Value())
		return nil, pkgerrors.NewInternalError("get children assets failed", err)
	}

	var assets []*asset.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.recordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *AssetRepository) FindByTypeAndGenre(ctx context.Context, assetType *asset.AssetType, genre *asset.Genre) ([]*asset.Asset, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (a:Asset)
		WHERE a.type = $type AND a.genre = $genre
		RETURN a
		ORDER BY a.createdAt DESC
	`

	params := map[string]interface{}{
		"type":  assetType.Value(),
		"genre": genre.Value(),
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to get assets by type and genre from Neo4j")
		return nil, pkgerrors.NewInternalError("get assets by type and genre failed", err)
	}

	var assets []*asset.Asset
	for result.Next() {
		record := result.Record()
		asset, err := r.recordToAsset(record)
		if err != nil {
			log.WithError(err).Error("Failed to convert record to asset")
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *AssetRepository) createParentRelationship(session neo4j.Session, childID, parentID string) error {
	query := `
		MATCH (child:Asset {id: $childId})
		MATCH (parent:Asset {id: $parentId})
		MERGE (child)-[:BELONGS_TO]->(parent)
		RETURN child, parent
	`

	params := map[string]interface{}{
		"childId":  childID,
		"parentId": parentID,
	}

	_, err := session.Run(query, params)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create parent relationship", "child_id", childID, "parent_id", parentID)
		return pkgerrors.NewInternalError("failed to create parent relationship", err)
	}

	return nil
}

func (r *AssetRepository) assetToParams(a *asset.Asset) map[string]interface{} {
	var title, description, assetType, genre, ownerID string
	if a.Title() != nil {
		title = a.Title().Value()
	}
	if a.Description() != nil {
		description = a.Description().Value()
	}
	if a.Type() != nil {
		assetType = a.Type().Value()
	}
	if a.Genre() != nil {
		genre = a.Genre().Value()
	}
	if a.OwnerID() != nil {
		ownerID = a.OwnerID().Value()
	}

	params := map[string]interface{}{
		"id":          a.ID().Value(),
		"slug":        a.Slug().Value(),
		"title":       title,
		"description": description,
		"type":        assetType,
		"genre":       genre,
		"genres":      r.genresToStringSlice(a.Genres()),
		"tags":        r.tagsToStringSlice(a.Tags()),
		"ownerId":     ownerID,
		"createdAt":   a.CreatedAt().Value().Format(time.RFC3339),
		"updatedAt":   a.UpdatedAt().Value().Format(time.RFC3339),
		"publishRule": nil,
	}

	if a.PublishRule() != nil {
		publishRuleData := map[string]interface{}{
			"publishAt":   a.PublishRule().PublishAt(),
			"unpublishAt": a.PublishRule().UnpublishAt(),
			"regions":     a.PublishRule().Regions(),
			"ageRating":   a.PublishRule().AgeRating(),
		}
		publishRuleJSON, _ := json.Marshal(publishRuleData)
		params["publishRule"] = string(publishRuleJSON)
	}

	var videosData []map[string]interface{}
	for _, video := range a.Videos() {
		videoData := map[string]interface{}{
			"id":     video.ID().Value(),
			"label":  video.Label().Value(),
			"type":   string(video.Type()),
			"format": string(video.Format()),
			"storageLocation": map[string]interface{}{
				"bucket": video.StorageLocation().Bucket(),
				"key":    video.StorageLocation().Key(),
				"url":    video.StorageLocation().URL(),
			},
			"width":              video.Width(),
			"height":             video.Height(),
			"duration":           video.Duration(),
			"bitrate":            video.Bitrate(),
			"codec":              video.Codec(),
			"size":               video.Size(),
			"contentType":        video.ContentType().Value(),
			"status":             string(video.Status()),
			"createdAt":          video.CreatedAt(),
			"updatedAt":          video.UpdatedAt(),
			"segmentCount":       video.SegmentCount(),
			"videoCodec":         video.VideoCodec(),
			"audioCodec":         video.AudioCodec(),
			"avgSegmentDuration": video.AvgSegmentDuration(),
			"segments":           video.Segments(),
			"frameRate":          video.FrameRate(),
			"audioChannels":      video.AudioChannels(),
			"audioSampleRate":    video.AudioSampleRate(),
		}
		if si := video.StreamInfo(); si != nil {
			videoData["streamInfo"] = map[string]interface{}{
				"downloadURL": si.DownloadURL(),
				"cdnPrefix":   si.CDNPrefix(),
				"url":         si.URL(),
			}
		}
		videosData = append(videosData, videoData)
	}
	videosJSON, _ := json.Marshal(videosData)

	var imagesData []map[string]interface{}
	for _, image := range a.Images() {
		imageData := map[string]interface{}{
			"id":          image.ID().Value(),
			"fileName":    image.FileName().Value(),
			"url":         image.URL().Value(),
			"type":        string(image.Type()),
			"contentType": image.ContentType().Value(),
			"createdAt":   image.CreatedAt().Format(time.RFC3339),
			"updatedAt":   image.UpdatedAt().Format(time.RFC3339),
		}
		if image.StorageLocation() != nil {
			imageData["storageLocation"] = map[string]interface{}{
				"bucket": image.StorageLocation().Bucket(),
				"key":    image.StorageLocation().Key(),
				"url":    image.StorageLocation().URL(),
			}
		}
		if image.Width() != nil {
			imageData["width"] = *image.Width()
		}
		if image.Height() != nil {
			imageData["height"] = *image.Height()
		}
		if image.Size() != nil {
			imageData["size"] = *image.Size()
		}
		if image.StreamInfo() != nil {
			imageData["streamInfo"] = map[string]interface{}{
				"downloadURL": image.StreamInfo().DownloadURL(),
				"cdnPrefix":   image.StreamInfo().CDNPrefix(),
				"url":         image.StreamInfo().URL(),
			}
		}
		if image.Metadata() != nil {
			imageData["metadata"] = image.Metadata()
		}
		imagesData = append(imagesData, imageData)
	}
	imagesJSON, _ := json.Marshal(imagesData)
	creditsJSON, _ := json.Marshal(a.Credits())
	metadataJSON, _ := json.Marshal(a.Metadata())

	params["videos"] = string(videosJSON)
	params["images"] = string(imagesJSON)
	params["credits"] = string(creditsJSON)
	params["metadata"] = string(metadataJSON)

	return params
}

func (r *AssetRepository) recordToAsset(record *neo4j.Record) (*asset.Asset, error) {
	log := r.logger

	assetNode, exists := record.Get("a")
	if !exists {
		return nil, pkgerrors.NewInternalError("asset node not found in record", nil)
	}

	node := assetNode.(neo4j.Node)
	props := node.Props

	id, _ := props["id"].(string)
	slug, _ := props["slug"].(string)
	title, _ := props["title"].(string)
	description, _ := props["description"].(string)
	assetType, _ := props["type"].(string)
	genre, _ := props["genre"].(string)
	ownerID, _ := props["ownerId"].(string)
	createdAtStr, _ := props["createdAt"].(string)
	updatedAtStr, _ := props["updatedAt"].(string)

	var genres []string
	if genresInterface, exists := props["genres"]; exists {
		if genresSlice, ok := genresInterface.([]interface{}); ok {
			for _, g := range genresSlice {
				if genreStr, ok := g.(string); ok {
					genres = append(genres, genreStr)
				}
			}
		}
	}

	var tags []string
	if tagsInterface, exists := props["tags"]; exists {
		if tagsSlice, ok := tagsInterface.([]interface{}); ok {
			for _, t := range tagsSlice {
				if tagStr, ok := t.(string); ok {
					tags = append(tags, tagStr)
				}
			}
		}
	}

	var publishRule *asset.PublishRule
	if publishRuleJSON, exists := props["publishRule"]; exists {
		if publishRuleStr, ok := publishRuleJSON.(string); ok && publishRuleStr != "" {
			var ruleData map[string]interface{}
			if err := json.Unmarshal([]byte(publishRuleStr), &ruleData); err != nil {
				log.WithError(err).Error("Failed to unmarshal publish rule JSON")
			} else {
				var publishAt *time.Time
				var unpublishAt *time.Time
				var regions []string
				var ageRating *string

				if publishAtVal, exists := ruleData["publishAt"]; exists && publishAtVal != nil {
					if publishAtStr, ok := publishAtVal.(string); ok {
						if t, err := time.Parse(time.RFC3339, publishAtStr); err == nil {
							publishAt = &t
						}
					}
				}

				if unpublishAtVal, exists := ruleData["unpublishAt"]; exists && unpublishAtVal != nil {
					if unpublishAtStr, ok := unpublishAtVal.(string); ok {
						if t, err := time.Parse(time.RFC3339, unpublishAtStr); err == nil {
							unpublishAt = &t
						}
					}
				}

				if regionsVal, exists := ruleData["regions"]; exists {
					if regionsSlice, ok := regionsVal.([]interface{}); ok {
						for _, r := range regionsSlice {
							if regionStr, ok := r.(string); ok {
								regions = append(regions, regionStr)
							}
						}
					}
				}

				if ageRatingVal, exists := ruleData["ageRating"]; exists && ageRatingVal != nil {
					if ageRatingStr, ok := ageRatingVal.(string); ok {
						ageRating = &ageRatingStr
					}
				}

				if publishAt != nil || unpublishAt != nil || len(regions) > 0 || ageRating != nil {
					publishRule, _ = asset.NewPublishRule(publishAt, unpublishAt, regions, ageRating)
				}
			}
		}
	}

	var videos []*asset.Video
	if videosJSON, exists := props["videos"]; exists {
		if videosStr, ok := videosJSON.(string); ok && videosStr != "" {

			var videosData []map[string]interface{}
			if err := json.Unmarshal([]byte(videosStr), &videosData); err != nil {
				log.WithError(err).Error("Failed to unmarshal videos JSON")
			} else {
				for _, videoData := range videosData {
					videoID, _ := videoData["id"].(string)
					label, _ := videoData["label"].(string)
					typeStr, _ := videoData["type"].(string)
					statusStr, _ := videoData["status"].(string)
					formatStr, _ := videoData["format"].(string)

					storageLocationMap, _ := videoData["storageLocation"].(map[string]interface{})
					bucket, _ := storageLocationMap["bucket"].(string)
					key, _ := storageLocationMap["key"].(string)
					url, _ := storageLocationMap["url"].(string)

					storageLocation, err := asset.NewS3Object(bucket, key, url)
					if err != nil {
						log.WithError(err).Error("Failed to create S3Object for video")
						continue
					}

					videoType := asset.VideoType(constants.VideoTypeMain)
					if typeStr != "" {
						videoType = asset.VideoType(typeStr)
					}

					videoFormat := asset.VideoFormat("")
					if formatStr != "" {
						videoFormat = asset.VideoFormat(formatStr)
					}

					videoStatus := asset.VideoStatus(constants.VideoStatusPending)
					if statusStr != "" {
						videoStatus = asset.VideoStatus(statusStr)
					}

					width, _ := videoData["width"].(float64)
					height, _ := videoData["height"].(float64)
					duration, _ := videoData["duration"].(float64)
					bitrate, _ := videoData["bitrate"].(float64)
					codec, _ := videoData["codec"].(string)
					size, _ := videoData["size"].(float64)
					contentType, _ := videoData["contentType"].(string)
					segmentCount, _ := videoData["segmentCount"].(float64)
					videoCodec, _ := videoData["videoCodec"].(string)
					audioCodec, _ := videoData["audioCodec"].(string)
					avgSegmentDuration, _ := videoData["avgSegmentDuration"].(float64)
					frameRate, _ := videoData["frameRate"].(string)
					audioChannels, _ := videoData["audioChannels"].(float64)
					audioSampleRate, _ := videoData["audioSampleRate"].(float64)

					var segments []string
					if segmentsInterface, ok := videoData["segments"].([]interface{}); ok {
						for _, seg := range segmentsInterface {
							if segStr, ok := seg.(string); ok {
								segments = append(segments, segStr)
							}
						}
					}

					createdAt, _ := videoData["createdAt"].(string)
					updatedAt, _ := videoData["updatedAt"].(string)

					var createdAtTime, updatedAtTime time.Time
					if createdAt != "" {
						if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
							createdAtTime = t
						} else {
							createdAtTime = time.Now().UTC()
						}
					} else {
						createdAtTime = time.Now().UTC()
					}

					if updatedAt != "" {
						if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
							updatedAtTime = t
						} else {
							updatedAtTime = time.Now().UTC()
						}
					} else {
						updatedAtTime = time.Now().UTC()
					}

					video, err := asset.ReconstructVideo(
						videoID,
						label,
						videoType,
						videoFormat,
						*storageLocation,
						int(width),
						int(height),
						duration,
						int(bitrate),
						codec,
						int64(size),
						contentType,
						videoStatus,
						createdAtTime,
						updatedAtTime,
						int(segmentCount),
						videoCodec,
						audioCodec,
						avgSegmentDuration,
						segments,
						frameRate,
						int(audioChannels),
						int(audioSampleRate),
					)
					if err != nil {
						log.WithError(err).Error("Failed to reconstruct video")
						continue
					}

					if streamInfoMap, ok := videoData["streamInfo"].(map[string]interface{}); ok {
						var downloadURL, cdnPrefix, urlStr *string
						if v, ok := streamInfoMap["downloadURL"].(string); ok && v != "" {
							downloadURL = &v
						}
						if v, ok := streamInfoMap["cdnPrefix"].(string); ok && v != "" {
							cdnPrefix = &v
						}
						if v, ok := streamInfoMap["url"].(string); ok && v != "" {
							urlStr = &v
						}
						if downloadURL != nil || cdnPrefix != nil || urlStr != nil {
							if si, err := asset.NewStreamInfo(downloadURL, cdnPrefix, urlStr); err == nil {
								video.SetStreamInfo(si)
							}
						}
					}

					videos = append(videos, video)
				}
			}
		}
	}

	var images []asset.Image
	if imagesJSON, exists := props["images"]; exists {
		if imagesStr, ok := imagesJSON.(string); ok && imagesStr != "" {
			var imageData []map[string]interface{}
			if err := json.Unmarshal([]byte(imagesStr), &imageData); err != nil {
				log.WithError(err).Error("Failed to unmarshal images JSON")
			} else {
				for _, imgData := range imageData {
					if img, err := r.reconstructImageFromData(imgData); err != nil {
						log.WithError(err).Error("Failed to reconstruct image from data")
					} else {
						images = append(images, *img)
					}
				}
			}
		}
	}

	var credits []asset.Credit
	if creditsJSON, exists := props["credits"]; exists {
		if creditsStr, ok := creditsJSON.(string); ok && creditsStr != "" {
			if err := json.Unmarshal([]byte(creditsStr), &credits); err != nil {
				log.WithError(err).Error("Failed to unmarshal credits JSON")
			}
		}
	}

	var metadata map[string]interface{}
	if metadataJSON, exists := props["metadata"]; exists {
		if metadataStr, ok := metadataJSON.(string); ok && metadataStr != "" {
			if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
				log.WithError(err).Error("Failed to unmarshal metadata JSON")
			}
		}
	}

	videosMap := make(map[string]*asset.Video)
	for _, video := range videos {
		videosMap[video.ID().Value()] = video
	}

	assetID, err := asset.NewAssetID(id)
	if err != nil {
		return nil, err
	}

	slugVO, err := asset.NewSlug(slug)
	if err != nil {
		return nil, err
	}

	var titleVO *asset.Title
	if title != "" {
		titleVO, err = asset.NewTitle(title)
		if err != nil {
			return nil, err
		}
	}

	var descriptionVO *asset.Description
	if description != "" {
		descriptionVO, err = asset.NewDescription(description)
		if err != nil {
			return nil, err
		}
	}

	var assetTypeVO *asset.AssetType
	if assetType != "" {
		assetTypeVO, err = asset.NewAssetType(assetType)
		if err != nil {
			return nil, err
		}
	}

	var genreVO *asset.Genre
	if genre != "" {
		genreVO, err = asset.NewGenre(genre)
		if err != nil {
			return nil, err
		}
	}

	genresVO, err := asset.NewGenres(genres)
	if err != nil {
		return nil, err
	}

	tagsVO, err := asset.NewTags(tags)
	if err != nil {
		return nil, err
	}

	var ownerIDVO *asset.OwnerID
	if ownerID != "" {
		ownerIDVO, err = asset.NewOwnerID(ownerID)
		if err != nil {
			return nil, err
		}
	}

	createdAtVO := ParseTimeToVO(createdAtStr, func(t time.Time) interface{} { return asset.NewCreatedAt(t) }).(*asset.CreatedAt)
	updatedAtVO := ParseTimeToVO(updatedAtStr, func(t time.Time) interface{} { return asset.NewUpdatedAt(t) }).(*asset.UpdatedAt)

	a := asset.ReconstructAsset(
		*assetID,
		*slugVO,
		titleVO,
		descriptionVO,
		assetTypeVO,
		genreVO,
		genresVO,
		tagsVO,
		*createdAtVO,
		*updatedAtVO,
		ownerIDVO,
		nil, //TODO: parentID will be set separately if needed
		images,
		videosMap,
		credits,
		publishRule,
		metadata,
	)

	return a, nil
}

func (r *AssetRepository) genresToStringSlice(genres *asset.Genres) []string {
	if genres == nil {
		return []string{}
	}

	genreStrings := make([]string, len(genres.Values()))
	for i, genre := range genres.Values() {
		genreStrings[i] = genre.Value()
	}
	return genreStrings
}

func (r *AssetRepository) tagsToStringSlice(tags *asset.Tags) []string {
	if tags == nil {
		return []string{}
	}

	tagStrings := make([]string, len(tags.Values()))
	for i, tag := range tags.Values() {
		tagStrings[i] = tag.Value()
	}
	return tagStrings
}

func (r *AssetRepository) reconstructImageFromData(imgData map[string]interface{}) (*asset.Image, error) {
	id, _ := imgData["id"].(string)
	fileName, _ := imgData["fileName"].(string)
	url, _ := imgData["url"].(string)
	typeStr, _ := imgData["type"].(string)

	var imageType asset.ImageType
	if typeStr != "" {
		imageType = asset.ImageType(typeStr)
	}

	var storageLocation *asset.S3Object
	if storageLocData, ok := imgData["storageLocation"].(map[string]interface{}); ok {
		bucket, _ := storageLocData["bucket"].(string)
		key, _ := storageLocData["key"].(string)
		urlStr, _ := storageLocData["url"].(string)
		if bucket != "" && key != "" && urlStr != "" {
			if s3Obj, err := asset.NewS3Object(bucket, key, urlStr); err == nil {
				storageLocation = s3Obj
			}
		}
	}

	var width *int
	if w, ok := imgData["width"].(float64); ok {
		widthInt := int(w)
		width = &widthInt
	}

	var height *int
	if h, ok := imgData["height"].(float64); ok {
		heightInt := int(h)
		height = &heightInt
	}

	var size *int64
	if s, ok := imgData["size"].(float64); ok {
		sizeInt64 := int64(s)
		size = &sizeInt64
	}

	var contentType string
	if ct, ok := imgData["contentType"].(string); ok {
		contentType = ct
	}

	var streamInfo *asset.StreamInfo
	if streamInfoData, ok := imgData["streamInfo"].(map[string]interface{}); ok {
		var downloadURL, cdnPrefix, urlStr *string
		if v, ok := streamInfoData["downloadURL"].(string); ok && v != "" {
			downloadURL = &v
		}
		if v, ok := streamInfoData["cdnPrefix"].(string); ok && v != "" {
			cdnPrefix = &v
		}
		if v, ok := streamInfoData["url"].(string); ok && v != "" {
			urlStr = &v
		}
		if downloadURL != nil || cdnPrefix != nil || urlStr != nil {
			if si, err := asset.NewStreamInfo(downloadURL, cdnPrefix, urlStr); err == nil {
				streamInfo = si
			}
		}
	}

	metadata := make(map[string]string)
	if metadataData, ok := imgData["metadata"].(map[string]interface{}); ok {
		for k, v := range metadataData {
			if str, ok := v.(string); ok {
				metadata[k] = str
			}
		}
	}

	createdAtStr, _ := imgData["createdAt"].(string)
	updatedAtStr, _ := imgData["updatedAt"].(string)

	var createdAt, updatedAt time.Time
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			createdAt = t
		} else {
			createdAt = time.Now().UTC()
		}
	} else {
		createdAt = time.Now().UTC()
	}

	if updatedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			updatedAt = t
		} else {
			updatedAt = time.Now().UTC()
		}
	} else {
		updatedAt = time.Now().UTC()
	}

	return asset.ReconstructImage(
		id,
		fileName,
		url,
		imageType,
		storageLocation,
		width,
		height,
		size,
		contentType,
		streamInfo,
		metadata,
		createdAt,
		updatedAt,
	)
}
