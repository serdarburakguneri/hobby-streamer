package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type Client struct {
	client *redis.Client
	logger *logger.Logger
}

type Service struct {
	client *Client
	ttl    TTLConfig
}

type TTLConfig struct {
	Bucket      time.Duration
	BucketsList time.Duration
	Asset       time.Duration
	AssetsList  time.Duration
}

type GraphQLBucketResponse struct {
	ID          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	AssetIDs    []string               `json:"assetIds"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	Assets      []GraphQLAssetResponse `json:"assets"`
}

type GraphQLAssetResponse struct {
	ID          string                     `json:"id"`
	Slug        string                     `json:"slug"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Type        string                     `json:"type"`
	Genre       string                     `json:"genre"`
	Genres      []string                   `json:"genres"`
	Tags        []string                   `json:"tags"`
	Status      string                     `json:"status"`
	Metadata    string                     `json:"metadata"`
	OwnerID     string                     `json:"ownerId"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
	Videos      []GraphQLVideoResponse     `json:"videos"`
	Images      []GraphQLImageResponse     `json:"images"`
	PublishRule GraphQLPublishRuleResponse `json:"publishRule"`
}

type GraphQLVideoResponse struct {
	ID                 string                    `json:"id"`
	Type               string                    `json:"type"`
	Format             string                    `json:"format"`
	StorageLocation    GraphQLS3ObjectResponse   `json:"storageLocation"`
	Width              int                       `json:"width"`
	Height             int                       `json:"height"`
	Duration           float64                   `json:"duration"`
	Bitrate            int                       `json:"bitrate"`
	Codec              string                    `json:"codec"`
	Size               int                       `json:"size"`
	ContentType        string                    `json:"contentType"`
	StreamInfo         GraphQLStreamInfoResponse `json:"streamInfo"`
	Status             string                    `json:"status"`
	CreatedAt          time.Time                 `json:"createdAt"`
	UpdatedAt          time.Time                 `json:"updatedAt"`
	IsReady            bool                      `json:"isReady"`
	IsProcessing       bool                      `json:"isProcessing"`
	IsFailed           bool                      `json:"isFailed"`
	SegmentCount       int                       `json:"segmentCount"`
	VideoCodec         string                    `json:"videoCodec"`
	AudioCodec         string                    `json:"audioCodec"`
	AvgSegmentDuration float64                   `json:"avgSegmentDuration"`
	Segments           []string                  `json:"segments"`
	FrameRate          string                    `json:"frameRate"`
	AudioChannels      int                       `json:"audioChannels"`
	AudioSampleRate    int                       `json:"audioSampleRate"`
}

type GraphQLImageResponse struct {
	ID              string                  `json:"id"`
	FileName        string                  `json:"fileName"`
	URL             string                  `json:"url"`
	Type            string                  `json:"type"`
	StorageLocation GraphQLS3ObjectResponse `json:"storageLocation"`
	Width           int                     `json:"width"`
	Height          int                     `json:"height"`
	Size            int                     `json:"size"`
	ContentType     string                  `json:"contentType"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

type GraphQLS3ObjectResponse struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type GraphQLStreamInfoResponse struct {
	DownloadURL string `json:"downloadUrl"`
	CDNPrefix   string `json:"cdnPrefix"`
	URL         string `json:"url"`
}

type GraphQLPublishRuleResponse struct {
	PublishAt   *time.Time `json:"publishAt"`
	UnpublishAt *time.Time `json:"unpublishAt"`
	Regions     []string   `json:"regions"`
	AgeRating   string     `json:"ageRating"`
}

func NewRedisClientWithConfig(host string, port int, db int, password string) (*Client, error) {
	if host == "" {
		return nil, errors.NewInternalError("Redis host is required", nil)
	}
	if port <= 0 {
		return nil, errors.NewInternalError("Redis port must be greater than 0", nil)
	}

	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		DB:           db,
		Password:     password,
		PoolSize:     20,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.NewTransientError("failed to connect to Redis", err)
	}

	return &Client{
		client: client,
		logger: logger.Get().WithService("redis-cache"),
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func NewService(client *Client, ttl TTLConfig) *Service {
	return &Service{client: client, ttl: ttl}
}

func (s *Service) GetBucket(ctx context.Context, key string) (*bucket.Bucket, error) {
	cacheKey := fmt.Sprintf("bucket:%s", key)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get bucket from cache", err)
	}

	var graphQLBucket GraphQLBucketResponse
	if err := json.Unmarshal(data, &graphQLBucket); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal bucket", err)
	}

	return s.convertGraphQLToDomainBucket(&graphQLBucket)
}

func (s *Service) SetBucket(ctx context.Context, bucket *bucket.Bucket) error {
	if bucket == nil {
		s.client.logger.Warn("Attempted to cache nil bucket, skipping")
		return nil
	}

	graphQLBucket := s.convertDomainToGraphQLBucket(bucket)
	cacheKey := fmt.Sprintf("bucket:%s", bucket.Key())

	data, err := json.Marshal(graphQLBucket)
	if err != nil {
		return errors.NewInternalError("failed to marshal bucket", err)
	}

	if err := s.client.client.Set(ctx, cacheKey, data, s.ttl.Bucket).Err(); err != nil {
		return errors.NewTransientError("failed to set bucket in cache", err)
	}

	return nil
}

func (s *Service) GetBuckets(ctx context.Context, limit int, nextKey *string) ([]*bucket.Bucket, error) {
	cacheKey := "buckets:list"

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get buckets from cache", err)
	}

	var graphQLBuckets []GraphQLBucketResponse
	if err := json.Unmarshal(data, &graphQLBuckets); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal buckets", err)
	}

	var buckets []*bucket.Bucket
	for _, graphQLBucket := range graphQLBuckets {
		domainBucket, err := s.convertGraphQLToDomainBucket(&graphQLBucket)
		if err != nil {
			s.client.logger.WithError(err).Warn("Failed to convert cached bucket to domain")
			continue
		}
		buckets = append(buckets, domainBucket)
	}

	start := 0
	if nextKey != nil {
		for i, b := range buckets {
			if b.ID().Value() == *nextKey {
				start = i + 1
				break
			}
		}
	}
	end := start + limit
	if end > len(buckets) {
		end = len(buckets)
	}
	return buckets[start:end], nil
}

func (s *Service) SetBuckets(ctx context.Context, buckets []*bucket.Bucket) error {
	cacheKey := "buckets:list"

	var graphQLBuckets []GraphQLBucketResponse
	for _, b := range buckets {
		if b == nil {
			s.client.logger.Warn("Nil bucket found before caching, skipping")
			continue
		}
		graphQLBucket := s.convertDomainToGraphQLBucket(b)
		graphQLBuckets = append(graphQLBuckets, *graphQLBucket)
	}

	if len(graphQLBuckets) == 0 {
		s.client.logger.Warn("No buckets to cache, skipping SetBuckets")
		return nil
	}

	data, err := json.Marshal(graphQLBuckets)
	if err != nil {
		return errors.NewInternalError("failed to marshal buckets", err)
	}

	if err := s.client.client.Set(ctx, cacheKey, data, s.ttl.BucketsList).Err(); err != nil {
		return errors.NewTransientError("failed to set buckets in cache", err)
	}

	return nil
}

func (s *Service) convertDomainToGraphQLBucket(b *bucket.Bucket) *GraphQLBucketResponse {
	graphQLBucket := &GraphQLBucketResponse{
		ID:        b.ID().Value(),
		Key:       b.Key().Value(),
		Name:      b.Name().Value(),
		Type:      b.Type().Value(),
		CreatedAt: b.CreatedAt().Value(),
		UpdatedAt: b.UpdatedAt().Value(),
	}

	if b.Description() != nil {
		graphQLBucket.Description = b.Description().Value()
	}
	if b.Status() != nil {
		graphQLBucket.Status = b.Status().Value()
	}
	if b.AssetIDs() != nil {
		graphQLBucket.AssetIDs = b.AssetIDs().Values()
	}

	for _, asset := range b.Assets() {
		graphQLAsset := s.convertDomainToGraphQLAsset(asset)
		graphQLBucket.Assets = append(graphQLBucket.Assets, *graphQLAsset)
	}

	return graphQLBucket
}

func (s *Service) convertGraphQLToDomainBucket(g *GraphQLBucketResponse) (*bucket.Bucket, error) {
	bucketID, err := bucket.NewBucketID(g.ID)
	if err != nil {
		return nil, err
	}

	bucketKey, err := bucket.NewBucketKey(g.Key)
	if err != nil {
		return nil, err
	}

	bucketName, err := bucket.NewBucketName(g.Name)
	if err != nil {
		return nil, err
	}

	var description *bucket.BucketDescription
	if g.Description != "" {
		descVO, err := bucket.NewBucketDescription(g.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	bucketType, err := bucket.NewBucketType(g.Type)
	if err != nil {
		return nil, err
	}

	var status *bucket.BucketStatus
	if g.Status != "" {
		statusVO, err := bucket.NewBucketStatus(g.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := bucket.NewCreatedAt(g.CreatedAt)
	updatedAt := bucket.NewUpdatedAt(g.UpdatedAt)

	assetIDs, err := bucket.NewAssetIDs(g.AssetIDs)
	if err != nil {
		return nil, err
	}

	var assets []*asset.Asset
	for _, graphQLAsset := range g.Assets {
		domainAsset, err := s.convertGraphQLToDomainAsset(&graphQLAsset)
		if err != nil {
			s.client.logger.WithError(err).Warn("Failed to convert cached asset to domain")
			continue
		}
		assets = append(assets, domainAsset)
	}

	return bucket.NewBucket(
		*bucketID,
		*bucketKey,
		*bucketName,
		description,
		*bucketType,
		status,
		assetIDs,
		*createdAt,
		*updatedAt,
		assets,
	), nil
}

func (s *Service) convertDomainToGraphQLAsset(a *asset.Asset) *GraphQLAssetResponse {
	graphQLAsset := &GraphQLAssetResponse{
		ID:        a.ID().Value(),
		Slug:      a.Slug().Value(),
		Type:      a.Type().Value(),
		Status:    a.Status().Value(),
		CreatedAt: a.CreatedAt().Value(),
		UpdatedAt: a.UpdatedAt().Value(),
	}

	if a.Title() != nil {
		graphQLAsset.Title = a.Title().Value()
	}
	if a.Description() != nil {
		graphQLAsset.Description = a.Description().Value()
	}
	if a.Genre() != nil {
		graphQLAsset.Genre = a.Genre().Value()
	}
	if a.Genres() != nil {
		genreValues := a.Genres().Values()
		graphQLAsset.Genres = make([]string, len(genreValues))
		for i, genre := range genreValues {
			graphQLAsset.Genres[i] = genre.Value()
		}
	}
	if a.Tags() != nil {
		graphQLAsset.Tags = a.Tags().Values()
	}
	if a.Metadata() != nil {
		graphQLAsset.Metadata = *a.Metadata()
	}
	if a.OwnerID() != nil {
		graphQLAsset.OwnerID = a.OwnerID().Value()
	}

	for _, video := range a.Videos() {
		graphQLVideo := s.convertDomainToGraphQLVideo(&video)
		graphQLAsset.Videos = append(graphQLAsset.Videos, *graphQLVideo)
	}

	for _, image := range a.Images() {
		graphQLImage := s.convertDomainToGraphQLImage(&image)
		graphQLAsset.Images = append(graphQLAsset.Images, *graphQLImage)
	}

	return graphQLAsset
}

func (s *Service) convertGraphQLToDomainAsset(g *GraphQLAssetResponse) (*asset.Asset, error) {
	assetID, err := asset.NewAssetID(g.ID)
	if err != nil {
		return nil, err
	}

	slug, err := asset.NewSlug(g.Slug)
	if err != nil {
		return nil, err
	}

	var title *asset.Title
	if g.Title != "" {
		titleVO, err := asset.NewTitle(g.Title)
		if err != nil {
			return nil, err
		}
		title = titleVO
	}

	var description *asset.Description
	if g.Description != "" {
		descVO, err := asset.NewDescription(g.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	assetType, err := asset.NewAssetType(g.Type)
	if err != nil {
		return nil, err
	}

	var genre *asset.Genre
	if g.Genre != "" {
		genreVO, err := asset.NewGenre(g.Genre)
		if err != nil {
			return nil, err
		}
		genre = genreVO
	}

	var genres *asset.Genres
	if len(g.Genres) > 0 {
		genresVO, err := asset.NewGenres(g.Genres)
		if err != nil {
			return nil, err
		}
		genres = genresVO
	}

	var tags *asset.Tags
	if len(g.Tags) > 0 {
		tagsVO, err := asset.NewTags(g.Tags)
		if err != nil {
			return nil, err
		}
		tags = tagsVO
	}

	var status *asset.Status
	if g.Status != "" {
		statusVO, err := asset.NewStatus(g.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := asset.NewCreatedAt(g.CreatedAt)
	updatedAt := asset.NewUpdatedAt(g.UpdatedAt)

	var ownerID *asset.OwnerID
	if g.OwnerID != "" {
		ownerIDVO, err := asset.NewOwnerID(g.OwnerID)
		if err != nil {
			return nil, err
		}
		ownerID = ownerIDVO
	}

	var videos []asset.Video
	for _, graphQLVideo := range g.Videos {
		domainVideo, err := s.convertGraphQLToDomainVideo(&graphQLVideo)
		if err != nil {
			s.client.logger.WithError(err).Warn("Failed to convert cached video to domain")
			continue
		}
		videos = append(videos, *domainVideo)
	}

	var images []asset.Image
	for _, graphQLImage := range g.Images {
		domainImage, err := s.convertGraphQLToDomainImage(&graphQLImage)
		if err != nil {
			s.client.logger.WithError(err).Warn("Failed to convert cached image to domain")
			continue
		}
		images = append(images, *domainImage)
	}

	return asset.NewAsset(
		*assetID,
		*slug,
		title,
		description,
		*assetType,
		genre,
		genres,
		tags,
		status,
		*createdAt,
		*updatedAt,
		&g.Metadata,
		ownerID,
		videos,
		images,
		nil,
	), nil
}

func (s *Service) convertDomainToGraphQLVideo(v *asset.Video) *GraphQLVideoResponse {
	graphQLVideo := &GraphQLVideoResponse{
		ID:                 v.ID().Value(),
		Format:             v.Format().Value(),
		Width:              *v.Width(),
		Height:             *v.Height(),
		Duration:           *v.Duration(),
		Bitrate:            *v.Bitrate(),
		Codec:              *v.Codec(),
		Size:               *v.Size(),
		ContentType:        *v.ContentType(),
		Status:             *v.Status(),
		CreatedAt:          v.CreatedAt(),
		UpdatedAt:          v.UpdatedAt(),
		IsReady:            v.IsReady(),
		IsProcessing:       v.IsProcessing(),
		IsFailed:           v.IsFailed(),
		SegmentCount:       *v.SegmentCount(),
		VideoCodec:         *v.VideoCodec(),
		AudioCodec:         *v.AudioCodec(),
		AvgSegmentDuration: *v.AvgSegmentDuration(),
		Segments:           v.Segments(),
		FrameRate:          *v.FrameRate(),
		AudioChannels:      *v.AudioChannels(),
		AudioSampleRate:    *v.AudioSampleRate(),
	}

	if v.Type() != nil {
		graphQLVideo.Type = v.Type().Value()
	}

	storageLocation := v.StorageLocation()
	graphQLVideo.StorageLocation = GraphQLS3ObjectResponse{
		Bucket: storageLocation.Bucket(),
		Key:    storageLocation.Key(),
		URL:    storageLocation.URL(),
	}

	if v.StreamInfo() != nil {
		streamInfo := v.StreamInfo()
		downloadURL := streamInfo.DownloadURL()
		cdnPrefix := streamInfo.CDNPrefix()
		url := streamInfo.URL()
		graphQLVideo.StreamInfo = GraphQLStreamInfoResponse{
			DownloadURL: *downloadURL,
			CDNPrefix:   *cdnPrefix,
			URL:         *url,
		}
	}

	return graphQLVideo
}

func (s *Service) convertGraphQLToDomainVideo(g *GraphQLVideoResponse) (*asset.Video, error) {
	videoID, err := asset.NewVideoID(g.ID)
	if err != nil {
		return nil, err
	}

	var videoType *asset.VideoType
	if g.Type != "" {
		videoTypeVO, err := asset.NewVideoType(g.Type)
		if err != nil {
			return nil, err
		}
		videoType = videoTypeVO
	}

	var format *asset.VideoFormat
	if g.Format != "" {
		formatVO, err := asset.NewVideoFormat(g.Format)
		if err != nil {
			return nil, err
		}
		format = formatVO
	}

	storageLocation, err := asset.NewS3ObjectValue(
		g.StorageLocation.Bucket,
		g.StorageLocation.Key,
		g.StorageLocation.URL,
	)
	if err != nil {
		return nil, err
	}

	var streamInfo *asset.StreamInfoValue
	if g.StreamInfo.DownloadURL != "" {
		streamInfoVO, err := asset.NewStreamInfoValue(
			&g.StreamInfo.DownloadURL,
			&g.StreamInfo.CDNPrefix,
			&g.StreamInfo.URL,
		)
		if err != nil {
			return nil, err
		}
		streamInfo = streamInfoVO
	}

	return asset.NewVideo(
		*videoID,
		videoType,
		format,
		*storageLocation,
		&g.Width,
		&g.Height,
		&g.Duration,
		&g.Bitrate,
		&g.Codec,
		&g.Size,
		&g.ContentType,
		streamInfo,
		nil,
		&g.Status,
		nil,
		g.CreatedAt,
		g.UpdatedAt,
		nil,
		g.IsReady,
		g.IsProcessing,
		g.IsFailed,
		&g.SegmentCount,
		&g.VideoCodec,
		&g.AudioCodec,
		&g.AvgSegmentDuration,
		g.Segments,
		&g.FrameRate,
		&g.AudioChannels,
		&g.AudioSampleRate,
		nil,
	), nil
}

func derefInt(p *int) int {
	if p != nil {
		return *p
	}
	return 0
}

func derefString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func (s *Service) convertDomainToGraphQLImage(i *asset.Image) *GraphQLImageResponse {
	graphQLImage := &GraphQLImageResponse{
		ID:          i.ID().Value(),
		FileName:    i.FileName().Value(),
		URL:         i.URL(),
		Width:       derefInt(i.Width()),
		Height:      derefInt(i.Height()),
		Size:        derefInt(i.Size()),
		ContentType: derefString(i.ContentType()),
		CreatedAt:   i.CreatedAt(),
		UpdatedAt:   i.UpdatedAt(),
	}

	if i.Type() != nil {
		graphQLImage.Type = i.Type().Value()
	}

	if i.StorageLocation() != nil {
		storageLocation := i.StorageLocation()
		graphQLImage.StorageLocation = GraphQLS3ObjectResponse{
			Bucket: storageLocation.Bucket(),
			Key:    storageLocation.Key(),
			URL:    storageLocation.URL(),
		}
	}

	return graphQLImage
}

func (s *Service) convertGraphQLToDomainImage(g *GraphQLImageResponse) (*asset.Image, error) {
	imageID, err := asset.NewImageID(g.ID)
	if err != nil {
		return nil, err
	}

	fileName, err := asset.NewFileName(g.FileName)
	if err != nil {
		return nil, err
	}

	var imageType *asset.ImageType
	if g.Type != "" {
		imageTypeVO, err := asset.NewImageType(g.Type)
		if err != nil {
			return nil, err
		}
		imageType = imageTypeVO
	}

	var storageLocation *asset.S3ObjectValue
	if g.StorageLocation.Bucket != "" {
		storageLocationVO, err := asset.NewS3ObjectValue(
			g.StorageLocation.Bucket,
			g.StorageLocation.Key,
			g.StorageLocation.URL,
		)
		if err != nil {
			return nil, err
		}
		storageLocation = storageLocationVO
	}

	return asset.NewImage(
		*imageID,
		*fileName,
		g.URL,
		imageType,
		storageLocation,
		&g.Width,
		&g.Height,
		&g.Size,
		&g.ContentType,
		nil,
		nil,
		g.CreatedAt,
		g.UpdatedAt,
	), nil
}

func (s *Service) GetAsset(ctx context.Context, slug string) (*asset.Asset, error) {
	cacheKey := fmt.Sprintf("asset:%s", slug)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get asset from cache", err)
	}

	var graphQLAsset GraphQLAssetResponse
	if err := json.Unmarshal(data, &graphQLAsset); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal asset", err)
	}

	return s.convertGraphQLToDomainAsset(&graphQLAsset)
}

func (s *Service) SetAsset(ctx context.Context, asset *asset.Asset) error {
	if asset == nil {
		s.client.logger.Warn("Attempted to cache nil asset, skipping")
		return nil
	}

	graphQLAsset := s.convertDomainToGraphQLAsset(asset)
	cacheKey := fmt.Sprintf("asset:%s", asset.Slug())

	data, err := json.Marshal(graphQLAsset)
	if err != nil {
		return errors.NewInternalError("failed to marshal asset", err)
	}

	if err := s.client.client.Set(ctx, cacheKey, data, s.ttl.Asset).Err(); err != nil {
		return errors.NewTransientError("failed to set asset in cache", err)
	}

	return nil
}

func (s *Service) GetAssets(ctx context.Context) ([]*asset.Asset, error) {
	cacheKey := "assets:list"

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get assets from cache", err)
	}

	var graphQLAssets []GraphQLAssetResponse
	if err := json.Unmarshal(data, &graphQLAssets); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal assets", err)
	}

	var assets []*asset.Asset
	for _, graphQLAsset := range graphQLAssets {
		domainAsset, err := s.convertGraphQLToDomainAsset(&graphQLAsset)
		if err != nil {
			s.client.logger.WithError(err).Warn("Failed to convert cached asset to domain")
			continue
		}
		assets = append(assets, domainAsset)
	}

	return assets, nil
}

func (s *Service) SetAssets(ctx context.Context, assets []*asset.Asset) error {
	cacheKey := "assets:list"

	var graphQLAssets []GraphQLAssetResponse
	for _, a := range assets {
		if a == nil {
			s.client.logger.Warn("Nil asset found before caching, skipping")
			continue
		}
		graphQLAsset := s.convertDomainToGraphQLAsset(a)
		graphQLAssets = append(graphQLAssets, *graphQLAsset)
	}

	data, err := json.Marshal(graphQLAssets)
	if err != nil {
		return errors.NewInternalError("failed to marshal assets", err)
	}

	if err := s.client.client.Set(ctx, cacheKey, data, s.ttl.AssetsList).Err(); err != nil {
		return errors.NewTransientError("failed to set assets in cache", err)
	}

	return nil
}

func (s *Service) InvalidateBucketCache(ctx context.Context, key string) error {
	cacheKey := fmt.Sprintf("bucket:%s", key)
	if err := s.client.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.NewTransientError("failed to invalidate bucket cache", err)
	}
	return nil
}

func (s *Service) InvalidateBucketsListCache(ctx context.Context) error {
	cacheKey := "buckets:list"
	if err := s.client.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.NewTransientError("failed to invalidate buckets list cache", err)
	}
	return nil
}

func (s *Service) InvalidateAssetCache(ctx context.Context, slug string) error {
	cacheKey := fmt.Sprintf("asset:%s", slug)
	if err := s.client.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.NewTransientError("failed to invalidate asset cache", err)
	}
	return nil
}

func (s *Service) InvalidateAssetsListCache(ctx context.Context) error {
	cacheKey := "assets:list"
	if err := s.client.client.Del(ctx, cacheKey).Err(); err != nil {
		return errors.NewTransientError("failed to invalidate assets list cache", err)
	}
	return nil
}
