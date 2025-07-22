package graphql

import (
	domainasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
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
		thumbMetadata := make([]string, 0, len(thumb.Metadata()))
		for k, v := range thumb.Metadata() {
			thumbMetadata = append(thumbMetadata, k+":"+v)
		}
		thumbnail = &Image{
			ID:          thumb.ID(),
			FileName:    thumb.FileName(),
			URL:         thumb.URL(),
			Type:        ImageType(thumb.Type()),
			Width:       thumb.Width(),
			Height:      thumb.Height(),
			Size:        &thumbSize,
			ContentType: thumb.ContentType(),
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
	contentType := video.ContentType()

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

	return &Video{
		ID:     video.ID(),
		Label:  video.Label(),
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

	return &Bucket{
		ID:          bucket.ID(),
		Key:         bucket.Key(),
		Name:        bucket.Name(),
		Description: &description,
		Type:        &bucketType,
		Status:      &status,
		OwnerID:     &ownerID,
		Metadata:    &metadata,
		CreatedAt:   bucket.CreatedAt(),
		UpdatedAt:   bucket.UpdatedAt(),
		AssetCount:  0, // Optionally, implement a method to count assets via relationships if needed
	}
}

func domainBucketPageToGraphQL(page *domainbucket.BucketPage) *BucketPage {
	if page == nil {
		return nil
	}

	items := make([]*Bucket, len(page.Items))
	for i, bucket := range page.Items {
		items[i] = domainBucketToGraphQL(bucket)
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

func parseMetadata(metadata *string) map[string]interface{} {
	if metadata == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["raw"] = *metadata
	return result
}
