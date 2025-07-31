package graphql

import (
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
)

func convertVideos(videos map[string]*assetentity.Video) []*Video {
	res := make([]*Video, 0, len(videos))
	for _, v := range videos {
		res = append(res, domainVideoToGraphQL(v))
	}
	return res
}

func convertImages(imgs []valueobjects.Image) []*Image {
	res := make([]*Image, len(imgs))
	for i, img := range imgs {
		res[i] = domainImageToGraphQL(&img)
	}
	return res
}

func domainAssetToGraphQL(asset *assetentity.Asset) *Asset {
	if asset == nil {
		return nil
	}

	title := ""
	if asset.Title() != nil {
		title = asset.Title().Value()
	}

	description := ""
	if asset.Description() != nil {
		description = asset.Description().Value()
	}

	assetType := ""
	if asset.Type() != nil {
		assetType = asset.Type().Value()
	}

	genre := ""
	if asset.Genre() != nil {
		genre = asset.Genre().Value()
	}

	ownerID := ""
	if asset.OwnerID() != nil {
		ownerID = asset.OwnerID().Value()
	}

	parentID := ""
	if asset.ParentID() != nil {
		parentID = asset.ParentID().Value()
	}

	genres := make([]string, 0)
	if asset.Genres() != nil {
		for _, genre := range asset.Genres().Values() {
			genres = append(genres, genre.Value())
		}
	}

	tags := make([]string, 0)
	if asset.Tags() != nil {
		for _, tag := range asset.Tags().Values() {
			tags = append(tags, tag.Value())
		}
	}

	metadata := ""
	if asset.Metadata() != nil {
		if raw, ok := asset.Metadata()["raw"]; ok {
			if str, ok := raw.(string); ok {
				metadata = str
			}
		}
	}

	var publishRule *PublishRule
	if asset.PublishRule() != nil {
		domainRule := asset.PublishRule()
		publishRule = &PublishRule{
			PublishAt:   domainRule.PublishAt(),
			UnpublishAt: domainRule.UnpublishAt(),
			Regions:     domainRule.Regions(),
			AgeRating:   domainRule.AgeRating(),
		}
	}

	videos := convertVideos(asset.Videos())

	images := convertImages(asset.Images())

	return &Asset{
		ID:          asset.ID().Value(),
		Slug:        asset.Slug().Value(),
		Title:       &title,
		Description: &description,
		Type:        &assetType,
		Genre:       &genre,
		OwnerID:     &ownerID,
		ParentID:    &parentID,
		Genres:      genres,
		Tags:        tags,
		Metadata:    &metadata,
		PublishRule: publishRule,
		Videos:      videos,
		Images:      images,
		CreatedAt:   asset.CreatedAt().Value(),
		UpdatedAt:   asset.UpdatedAt().Value(),
	}
}
func domainVideoToGraphQL(video *assetentity.Video) *Video {
	if video == nil {
		return nil
	}

	storageLocation := video.StorageLocation()
	var gqlStorageLocation *S3Object
	if storageLocation != (valueobjects.S3Object{}) {
		gqlStorageLocation = &S3Object{
			Bucket: storageLocation.Bucket(),
			Key:    storageLocation.Key(),
			URL:    storageLocation.URL(),
		}
	}

	var streamInfo *StreamInfo
	if video.StreamInfo() != nil {
		domainStreamInfo := video.StreamInfo()
		streamInfo = &StreamInfo{
			DownloadURL: domainStreamInfo.DownloadURL(),
			CdnPrefix:   domainStreamInfo.CDNPrefix(),
			URL:         domainStreamInfo.URL(),
		}
	}

	width := video.Width()
	height := video.Height()
	duration := video.Duration()
	bitrate := video.Bitrate()
	codec := video.Codec()
	size := int(video.Size())
	contentType := video.ContentType().Value()
	format := VideoFormat(video.Format().Value())
	status := VideoStatus(video.Status().Value())
	videoType := VideoType(video.Type().Value())
	segmentCount := video.SegmentCount()
	videoCodec := video.VideoCodec()
	audioCodec := video.AudioCodec()
	avgSegmentDuration := video.AvgSegmentDuration()
	frameRate := video.FrameRate()
	audioChannels := video.AudioChannels()
	audioSampleRate := video.AudioSampleRate()

	return &Video{
		ID:                 video.ID().Value(),
		Label:              video.Label().Value(),
		Type:               videoType,
		Format:             &format,
		StorageLocation:    gqlStorageLocation,
		Width:              &width,
		Height:             &height,
		Duration:           &duration,
		Bitrate:            &bitrate,
		Codec:              &codec,
		Size:               &size,
		ContentType:        &contentType,
		StreamInfo:         streamInfo,
		Status:             status,
		CreatedAt:          video.CreatedAt(),
		UpdatedAt:          video.UpdatedAt(),
		IsReady:            video.IsReady(),
		IsProcessing:       video.IsProcessing(),
		IsFailed:           video.IsFailed(),
		SegmentCount:       &segmentCount,
		VideoCodec:         &videoCodec,
		AudioCodec:         &audioCodec,
		AvgSegmentDuration: &avgSegmentDuration,
		Segments:           video.Segments(),
		FrameRate:          &frameRate,
		AudioChannels:      &audioChannels,
		AudioSampleRate:    &audioSampleRate,
	}
}
func domainImageToGraphQL(img *valueobjects.Image) *Image {
	if img == nil {
		return nil
	}

	var streamInfo *StreamInfo
	if img.StreamInfo() != nil {
		domainStreamInfo := img.StreamInfo()
		streamInfo = &StreamInfo{
			DownloadURL: domainStreamInfo.DownloadURL(),
			CdnPrefix:   domainStreamInfo.CDNPrefix(),
			URL:         domainStreamInfo.URL(),
		}
	}

	width := img.Width()
	height := img.Height()
	size := img.Size()
	contentType := img.ContentType().Value()
	var sizeInt *int
	if size != nil {
		sizeVal := int(*size)
		sizeInt = &sizeVal
	}

	var storageLocation *S3Object
	if img.StorageLocation() != nil {
		domainStorageLocation := img.StorageLocation()
		storageLocation = &S3Object{
			Bucket: domainStorageLocation.Bucket(),
			Key:    domainStorageLocation.Key(),
			URL:    domainStorageLocation.URL(),
		}
	}

	return &Image{
		ID:              img.ID().Value(),
		FileName:        img.FileName().Value(),
		URL:             img.URL(),
		Type:            ImageType(img.Type().Value()),
		StorageLocation: storageLocation,
		Width:           width,
		Height:          height,
		Size:            sizeInt,
		ContentType:     &contentType,
		StreamInfo:      streamInfo,
		Metadata:        []string{},
		CreatedAt:       img.CreatedAt(),
		UpdatedAt:       img.UpdatedAt(),
	}
}
func domainAssetPageToGraphQL(page *assetentity.AssetPage) *AssetPage {
	if page == nil {
		return nil
	}

	items := make([]*Asset, len(page.Items))
	for i, asset := range page.Items {
		items[i] = domainAssetToGraphQL(asset)
	}

	var nextKeyStr *string
	if page.LastKey != nil {
		if key, ok := page.LastKey["key"].(string); ok {
			nextKeyStr = &key
		}
	}

	return &AssetPage{
		Items:   items,
		NextKey: nextKeyStr,
		HasMore: page.HasMore,
	}
}
func domainBucketToGraphQL(bucket *bucketentity.Bucket) *Bucket {
	if bucket == nil {
		return nil
	}

	description := ""
	if bucket.Description() != nil {
		description = bucket.Description().Value()
	}

	ownerID := ""
	if bucket.OwnerID() != nil {
		ownerID = bucket.OwnerID().Value()
	}

	status := ""
	if bucket.Status() != nil {
		status = bucket.Status().Value()
	}

	bucketType := ""
	if bucket.Type() != nil {
		bucketType = bucket.Type().Value()
	}

	metadata := ""
	if bucket.Metadata() != nil {
		if raw, ok := bucket.Metadata()["raw"]; ok {
			if str, ok := raw.(string); ok {
				metadata = str
			}
		}
	}

	return &Bucket{
		ID:          bucket.ID().Value(),
		Key:         bucket.Key().Value(),
		Name:        bucket.Name().Value(),
		Description: &description,
		Type:        bucketType,
		Status:      &status,
		OwnerID:     &ownerID,
		Metadata:    &metadata,
		CreatedAt:   bucket.CreatedAt().Value(),
		UpdatedAt:   bucket.UpdatedAt().Value(),
	}
}
func domainBucketPageToGraphQL(page *bucketentity.BucketPage) *BucketPage {
	if page == nil {
		return nil
	}

	items := make([]*Bucket, len(page.Items))
	for i, bucket := range page.Items {
		items[i] = domainBucketToGraphQL(bucket)
	}

	var nextKeyStr *string
	if page.LastKey != nil {
		if key, ok := page.LastKey["key"].(string); ok {
			nextKeyStr = &key
		}
	}

	return &BucketPage{
		Items:   items,
		NextKey: nextKeyStr,
		HasMore: page.HasMore,
	}
}
