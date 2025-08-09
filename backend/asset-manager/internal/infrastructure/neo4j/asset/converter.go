package asset

import (
	"encoding/json"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AssetConverter struct {
	logger *logger.Logger
}

func NewAssetConverter(logger *logger.Logger) *AssetConverter {
	return &AssetConverter{
		logger: logger,
	}
}

func (c *AssetConverter) AssetToParams(a *entity.Asset) map[string]interface{} {
	var title, description, assetType, genre, ownerID, parentID string
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
	if a.ParentID() != nil {
		parentID = a.ParentID().Value()
	}

	params := map[string]interface{}{
		"id":              a.ID().Value(),
		"expectedVersion": a.Version(),
		"slug":            a.Slug().Value(),
		"title":           title,
		"description":     description,
		"type":            assetType,
		"genre":           genre,
		"genres":          c.genresToStringSlice(a.Genres()),
		"tags":            c.tagsToStringSlice(a.Tags()),
		"ownerId":         ownerID,
		"parentId":        parentID,
		"createdAt":       a.CreatedAt().Value().Format(time.RFC3339),
		"updatedAt":       a.UpdatedAt().Value().Format(time.RFC3339),
		"publishRule":     nil,
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
			"createdAt":          video.Timestamps().CreatedAt(),
			"updatedAt":          video.Timestamps().UpdatedAt(),
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
	params["videos"] = string(videosJSON)

	var imagesData []map[string]interface{}
	for _, image := range a.Images() {
		imageData := map[string]interface{}{
			"id":          image.ID().Value(),
			"fileName":    image.FileName().Value(),
			"url":         image.URL(),
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
	params["images"] = string(imagesJSON)

	creditsJSON, _ := json.Marshal(a.Credits())
	params["credits"] = string(creditsJSON)

	metadataJSON, _ := json.Marshal(a.Metadata())
	params["metadata"] = string(metadataJSON)

	return params
}

func (c *AssetConverter) RecordToAsset(record *neo4j.Record) (*entity.Asset, error) {
	log := c.logger

	assetNode, exists := record.Get("a")
	if !exists {
		return nil, pkgerrors.NewInternalError("asset node not found in record", nil)
	}

	node := assetNode.(neo4j.Node)
	props := node.Props

	id, _ := props["id"].(string)

	var version int
	if v, ok := props["version"].(int64); ok {
		version = int(v)
	}
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

	var publishRule *valueobjects.PublishRule
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
					publishRule, _ = valueobjects.NewPublishRule(publishAt, unpublishAt, regions, ageRating)
				}
			}
		}
	}

	var videos []*entity.Video
	if videosJSON, exists := props["videos"]; exists {
		if videosStr, ok := videosJSON.(string); ok && videosStr != "" {
			var videosData []map[string]interface{}
			if err := json.Unmarshal([]byte(videosStr), &videosData); err != nil {
				log.WithError(err).Error("Failed to unmarshal videos JSON")
			} else {
				for _, videoData := range videosData {
					if video, err := c.reconstructVideoFromData(videoData); err == nil {
						videos = append(videos, video)
					}
				}
			}
		}
	}

	var images []valueobjects.Image
	if imagesJSON, exists := props["images"]; exists {
		if imagesStr, ok := imagesJSON.(string); ok && imagesStr != "" {
			var imageData []map[string]interface{}
			if err := json.Unmarshal([]byte(imagesStr), &imageData); err != nil {
				log.WithError(err).Error("Failed to unmarshal images JSON")
			} else {
				for _, imgData := range imageData {
					if img, err := c.reconstructImageFromData(imgData); err != nil {
						log.WithError(err).Error("Failed to reconstruct image from data")
					} else {
						images = append(images, *img)
					}
				}
			}
		}
	}

	var credits []valueobjects.Credit
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

	videosMap := make(map[string]*entity.Video)
	for _, video := range videos {
		videosMap[video.ID().Value()] = video
	}

	assetID, err := valueobjects.NewAssetID(id)
	if err != nil {
		return nil, err
	}

	slugVO, err := valueobjects.NewSlug(slug)
	if err != nil {
		return nil, err
	}

	var titleVO *valueobjects.Title
	if title != "" {
		titleVO, err = valueobjects.NewTitle(title)
		if err != nil {
			return nil, err
		}
	}

	var descriptionVO *valueobjects.Description
	if description != "" {
		descriptionVO, err = valueobjects.NewDescription(description)
		if err != nil {
			return nil, err
		}
	}

	var assetTypeVO *valueobjects.AssetType
	if assetType != "" {
		assetTypeVO, err = valueobjects.NewAssetType(assetType)
		if err != nil {
			return nil, err
		}
	}

	var genreVO *valueobjects.Genre
	if genre != "" {
		genreVO, err = valueobjects.NewGenre(genre)
		if err != nil {
			return nil, err
		}
	}

	genresVO, err := valueobjects.NewGenres(genres)
	if err != nil {
		return nil, err
	}

	tagsVO, err := valueobjects.NewTags(tags)
	if err != nil {
		return nil, err
	}

	var ownerIDVO *valueobjects.OwnerID
	if ownerID != "" {
		ownerIDVO, err = valueobjects.NewOwnerID(ownerID)
		if err != nil {
			return nil, err
		}
	}

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to parse createdAt", err)
	}
	createdAtVO := valueobjects.NewCreatedAt(createdAtTime)

	updatedAtTime, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to parse updatedAt", err)
	}
	updatedAtVO := valueobjects.NewUpdatedAt(updatedAtTime)

	a := entity.ReconstructAsset(
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

	a.SetVersion(version)

	return a, nil
}

func (c *AssetConverter) reconstructImageFromData(imgData map[string]interface{}) (*valueobjects.Image, error) {
	idVO, err := valueobjects.NewID(imgData["id"].(string), "image id", 36)
	if err != nil {
		return nil, err
	}
	fileName, _ := imgData["fileName"].(string)
	url, _ := imgData["url"].(string)
	typeStr, _ := imgData["type"].(string)

	imageType, err := valueobjects.NewImageType(typeStr)
	if err != nil {
		return nil, err
	}

	var storageLocation *valueobjects.S3Object
	if storageLocData, ok := imgData["storageLocation"].(map[string]interface{}); ok {
		bucket, _ := storageLocData["bucket"].(string)
		key, _ := storageLocData["key"].(string)
		urlStr, _ := storageLocData["url"].(string)
		if bucket != "" && key != "" && urlStr != "" {
			if s3Obj, err := valueobjects.NewS3Object(bucket, key, urlStr); err == nil {
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

	var streamInfo *valueobjects.StreamInfo
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
			if si, err := valueobjects.NewStreamInfo(downloadURL, cdnPrefix, urlStr); err == nil {
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

	return valueobjects.NewImageWithDetails(
		*idVO,
		fileName,
		url,
		*imageType,
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

func (c *AssetConverter) reconstructVideoFromData(videoData map[string]interface{}) (*entity.Video, error) {
	log := c.logger
	videoID, err := valueobjects.NewVideoID(videoData["id"].(string))
	if err != nil {
		log.WithError(err).Error("Failed to create video ID")
		return nil, err
	}
	label, _ := videoData["label"].(string)
	typeStr, _ := videoData["type"].(string)
	statusStr, _ := videoData["status"].(string)
	formatStr, _ := videoData["format"].(string)

	storageLocationMap, _ := videoData["storageLocation"].(map[string]interface{})
	bucket, _ := storageLocationMap["bucket"].(string)
	key, _ := storageLocationMap["key"].(string)
	url, _ := storageLocationMap["url"].(string)

	storageLocation, err := valueobjects.NewS3Object(bucket, key, url)
	if err != nil {
		log.WithError(err).Error("Failed to create S3Object for video")
		return nil, err
	}

	videoType, err := valueobjects.NewVideoType(typeStr)
	if err != nil {
		log.WithError(err).Error("Failed to create video type")
		return nil, err
	}

	videoFormat, err := valueobjects.NewVideoFormat(formatStr)
	if err != nil {
		log.WithError(err).Error("Failed to create video format")
		return nil, err
	}

	videoStatus, err := valueobjects.NewVideoStatus(statusStr)
	if err != nil {
		log.WithError(err).Error("Failed to create video status")
		return nil, err
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

	timestamps := valueobjects.NewTimestamps()

	labelVO, err := valueobjects.NewValidatedString(label, 100, "label")
	if err != nil {
		return nil, err
	}
	contentTypeVO, err := valueobjects.NewContentType(contentType)
	if err != nil {
		return nil, err
	}

	video := entity.ReconstructVideo(
		*videoID,
		*labelVO,
		*videoType,
		*videoFormat,
		*storageLocation,
		int(width),
		int(height),
		duration,
		int(bitrate),
		codec,
		int64(size),
		*contentTypeVO,
		*videoStatus,
		timestamps,
		int(segmentCount),
		videoCodec,
		audioCodec,
		avgSegmentDuration,
		segments,
		frameRate,
		int(audioChannels),
		int(audioSampleRate),
	)

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
			if si, err := valueobjects.NewStreamInfo(downloadURL, cdnPrefix, urlStr); err == nil {
				video.SetStreamInfo(si)
			}
		}
	}
	return video, nil
}

func (c *AssetConverter) genresToStringSlice(genres *valueobjects.Genres) []string {
	if genres == nil {
		return []string{}
	}
	genreStrings := make([]string, len(genres.Values()))
	for i, genre := range genres.Values() {
		genreStrings[i] = genre.Value()
	}
	return genreStrings
}

func (c *AssetConverter) tagsToStringSlice(tags *valueobjects.Tags) []string {
	if tags == nil {
		return []string{}
	}
	tagStrings := make([]string, len(tags.Values()))
	for i, tag := range tags.Values() {
		tagStrings[i] = tag.Value()
	}
	return tagStrings
}
