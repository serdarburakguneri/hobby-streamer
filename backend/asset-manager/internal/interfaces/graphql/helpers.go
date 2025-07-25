package graphql

import (
	"fmt"

	domainasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func domainAssetToGraphQL(asset *domainasset.Asset) *Asset {
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

	videos := make([]*Video, 0, len(asset.Videos()))
	for _, video := range asset.Videos() {
		videos = append(videos, domainVideoToGraphQL(video))
	}

	images := make([]*Image, 0, len(asset.Images()))
	for _, image := range asset.Images() {
		images = append(images, domainImageToGraphQL(&image))
	}

	return &Asset{
		ID:          asset.ID().Value(),
		Slug:        asset.Slug().Value(),
		Title:       &title,
		Description: &description,
		Type:        &assetType,
		Genre:       &genre,
		Genres:      genres,
		Tags:        tags,
		CreatedAt:   asset.CreatedAt().Value(),
		UpdatedAt:   asset.UpdatedAt().Value(),
		OwnerID:     &ownerID,
		ParentID:    &parentID,
		Metadata:    &metadata,
		Status:      asset.Status(),
		PublishRule: publishRule,
		Videos:      videos,
		Images:      images,
	}
}

func domainVideoToGraphQL(video *domainasset.Video) *Video {
	if video == nil {
		return nil
	}

	var streamInfo *StreamInfo
	if video.StreamInfo() != nil {
		streamInfo = &StreamInfo{
			DownloadURL: video.StreamInfo().DownloadURL(),
			CdnPrefix:   video.StreamInfo().CDNPrefix(),
			URL:         video.StreamInfo().URL(),
		}
	}

	var transcodingInfo *TranscodingInfo
	ti := video.TranscodingInfo()
	if ti.JobID() != "" {
		jobID := ti.JobID()
		progress := ti.Progress()
		outputURL := ti.OutputURL()
		transcodingInfo = &TranscodingInfo{
			JobID:       &jobID,
			Progress:    &progress,
			OutputURL:   &outputURL,
			Error:       ti.Error(),
			CompletedAt: ti.CompletedAt(),
		}
	}

	var thumbnail *Image
	if video.Thumbnail() != nil {
		thumb := video.Thumbnail()
		thumbSize := int(*thumb.Size())
		thumbContentType := thumb.ContentType().Value()
		thumbMetadata := make([]string, 0, len(thumb.Metadata()))
		for k, v := range thumb.Metadata() {
			thumbMetadata = append(thumbMetadata, k+":"+v)
		}
		thumbnail = &Image{
			ID:          thumb.ID().Value(),
			FileName:    thumb.FileName().Value(),
			URL:         thumb.URL().Value(),
			Type:        ImageType(thumb.Type()),
			Width:       thumb.Width(),
			Height:      thumb.Height(),
			Size:        &thumbSize,
			ContentType: &thumbContentType,
			Metadata:    thumbMetadata,
			CreatedAt:   thumb.CreatedAt(),
			UpdatedAt:   thumb.UpdatedAt(),
		}
	}

	width := video.Width()
	height := video.Height()
	duration := video.Duration()
	bitrate := video.Bitrate()
	codec := video.Codec()
	size := int(video.Size())
	contentType := video.ContentType().Value()

	videoMetadata := make([]string, 0, len(video.Metadata()))
	for k, v := range video.Metadata() {
		videoMetadata = append(videoMetadata, k+":"+v)
	}

	format := video.Format()
	var formatPtr *VideoFormat
	if format != "" {
		formatStr := string(format)
		formatEnum := VideoFormat(formatStr)
		formatPtr = &formatEnum
	}

	segmentCount := video.SegmentCount()
	videoCodec := video.VideoCodec()
	audioCodec := video.AudioCodec()
	avgSegmentDuration := video.AvgSegmentDuration()
	segments := video.Segments()
	frameRate := video.FrameRate()
	audioChannels := video.AudioChannels()
	audioSampleRate := video.AudioSampleRate()

	return &Video{
		ID:     video.ID().Value(),
		Label:  video.Label().Value(),
		Type:   VideoType(video.Type()),
		Format: formatPtr,
		StorageLocation: &S3Object{
			Bucket: video.StorageLocation().Bucket(),
			Key:    video.StorageLocation().Key(),
			URL:    video.StorageLocation().URL(),
		},
		Width:              &width,
		Height:             &height,
		Duration:           &duration,
		Bitrate:            &bitrate,
		Codec:              &codec,
		Size:               &size,
		ContentType:        &contentType,
		StreamInfo:         streamInfo,
		Metadata:           videoMetadata,
		Status:             VideoStatus(video.Status()),
		Thumbnail:          thumbnail,
		TranscodingInfo:    transcodingInfo,
		CreatedAt:          video.CreatedAt(),
		UpdatedAt:          video.UpdatedAt(),
		Quality:            VideoQuality(video.Quality()),
		IsReady:            video.IsReady(),
		IsProcessing:       video.IsProcessing(),
		IsFailed:           video.IsFailed(),
		SegmentCount:       &segmentCount,
		VideoCodec:         &videoCodec,
		AudioCodec:         &audioCodec,
		AvgSegmentDuration: &avgSegmentDuration,
		Segments:           segments,
		FrameRate:          &frameRate,
		AudioChannels:      &audioChannels,
		AudioSampleRate:    &audioSampleRate,
	}
}

func domainBucketToGraphQL(bucket *domainbucket.Bucket) *Bucket {
	if bucket == nil {
		return nil
	}

	description := ""
	if bucket.Description() != nil {
		description = *bucket.Description()
	}

	ownerID := ""
	if bucket.OwnerID() != nil {
		ownerID = *bucket.OwnerID()
	}

	bucketType := ""
	if bucket.Type() != nil {
		bucketType = *bucket.Type()
	}

	status := ""
	if bucket.Status() != nil {
		status = *bucket.Status()
	}

	metadata := ""
	if bucket.Metadata() != nil {
		if raw, ok := bucket.Metadata()["raw"]; ok {
			if str, ok := raw.(string); ok {
				metadata = str
			}
		}
	}

	var assets []*Asset = nil

	return &Bucket{
		ID:          bucket.ID().Value(),
		Key:         bucket.Key(),
		Name:        bucket.Name(),
		Description: &description,
		Type:        bucketType,
		Status:      &status,
		OwnerID:     &ownerID,
		Assets:      assets,
		Metadata:    &metadata,
		CreatedAt:   bucket.CreatedAt(),
		UpdatedAt:   bucket.UpdatedAt(),
	}
}

func domainBucketPageToGraphQL(page *domainbucket.BucketPage) *BucketPage {
	if page == nil {
		return nil
	}

	items := make([]*Bucket, 0, len(page.Items))
	for i, bucket := range page.Items {
		if bucket == nil {
			logger.Get().Error(fmt.Sprintf("nil bucket encountered in BucketPage.Items at index %d", i))
			continue
		}
		items = append(items, domainBucketToGraphQL(bucket))
	}

	nextKey := ""
	if page.LastKey != nil {
		if key, ok := page.LastKey["key"]; ok {
			if str, ok := key.(string); ok {
				nextKey = str
			}
		}
	}

	return &BucketPage{
		Items:   items,
		NextKey: &nextKey,
		HasMore: page.HasMore,
	}
}

func domainAssetPageToGraphQL(page *domainasset.AssetPage) *AssetPage {
	if page == nil {
		return nil
	}

	items := make([]*Asset, len(page.Items))
	for i, asset := range page.Items {
		items[i] = domainAssetToGraphQL(asset)
	}

	nextKey := ""
	if page.LastKey != nil {
		if key, ok := page.LastKey["key"]; ok {
			if str, ok := key.(string); ok {
				nextKey = str
			}
		}
	}

	return &AssetPage{
		Items:   items,
		NextKey: &nextKey,
		HasMore: page.HasMore,
	}
}

func domainImageToGraphQL(img *domainasset.Image) *Image {
	if img == nil {
		return nil
	}

	var storageLocation *S3Object
	if img.StorageLocation() != nil {
		storageLocation = &S3Object{
			Bucket: img.StorageLocation().Bucket(),
			Key:    img.StorageLocation().Key(),
			URL:    img.StorageLocation().URL(),
		}
	}

	var width *int
	if img.Width() != nil {
		width = img.Width()
	}

	var height *int
	if img.Height() != nil {
		height = img.Height()
	}

	var size *int
	if img.Size() != nil {
		s := int(*img.Size())
		size = &s
	}

	var contentType *string
	contentTypeValue := img.ContentType().Value()
	if contentTypeValue != "" {
		contentType = &contentTypeValue
	}

	var streamInfo *StreamInfo
	if img.StreamInfo() != nil {
		streamInfo = &StreamInfo{
			DownloadURL: img.StreamInfo().DownloadURL(),
			CdnPrefix:   img.StreamInfo().CDNPrefix(),
			URL:         img.StreamInfo().URL(),
		}
	}

	metadata := make([]string, 0)
	if img.Metadata() != nil {
		for k, v := range img.Metadata() {
			metadata = append(metadata, k+":"+v)
		}
	}

	return &Image{
		ID:              img.ID().Value(),
		FileName:        img.FileName().Value(),
		URL:             img.URL().Value(),
		Type:            ImageType(img.Type()),
		StorageLocation: storageLocation,
		Width:           width,
		Height:          height,
		Size:            size,
		ContentType:     contentType,
		StreamInfo:      streamInfo,
		Metadata:        metadata,
		CreatedAt:       img.CreatedAt(),
		UpdatedAt:       img.UpdatedAt(),
	}
}

func parseMetadata(metadata *string) map[string]interface{} {
	if metadata == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["raw"] = *metadata
	return result
}
