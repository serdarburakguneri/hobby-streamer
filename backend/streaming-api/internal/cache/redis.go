package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/model"
)

type Client struct {
	client *redis.Client
	logger *logger.Logger
}

type Service struct {
	client *Client
}

func NewRedisClientWithConfig(host string, port int, db int, password string) (*Client, error) {
	if host == "" {
		return nil, errors.NewInternalError("Redis host is required", nil)
	}
	if port <= 0 {
		return nil, errors.NewInternalError("Redis port must be greater than 0", nil)
	}

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		DB:       db,
		Password: password,
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

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetBucket(ctx context.Context, key string) (*model.Bucket, error) {
	cacheKey := fmt.Sprintf("bucket:%s", key)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get bucket from cache", err)
	}

	var bucket model.Bucket
	if err := json.Unmarshal(data, &bucket); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal bucket", err)
	}

	return &bucket, nil
}

func (s *Service) SetBucket(ctx context.Context, bucket *model.Bucket) error {
	cacheKey := fmt.Sprintf("bucket:%s", bucket.Key)

	data, err := json.Marshal(bucket)
	if err != nil {
		return errors.NewInternalError("failed to marshal bucket", err)
	}

	ttl := 30 * time.Minute
	if err := s.client.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return errors.NewTransientError("failed to set bucket in cache", err)
	}

	return nil
}

func (s *Service) GetBuckets(ctx context.Context) ([]model.Bucket, error) {
	cacheKey := "buckets:list"

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get buckets from cache", err)
	}

	var buckets []model.Bucket
	if err := json.Unmarshal(data, &buckets); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal buckets", err)
	}

	return buckets, nil
}

func (s *Service) SetBuckets(ctx context.Context, buckets []model.Bucket) error {
	cacheKey := "buckets:list"

	data, err := json.Marshal(buckets)
	if err != nil {
		return errors.NewInternalError("failed to marshal buckets", err)
	}

	ttl := 30 * time.Minute
	if err := s.client.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return errors.NewTransientError("failed to set buckets in cache", err)
	}

	return nil
}

func (s *Service) GetAsset(ctx context.Context, slug string) (*model.Asset, error) {
	cacheKey := fmt.Sprintf("asset:%s", slug)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get asset from cache", err)
	}

	var asset model.Asset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal asset", err)
	}

	return &asset, nil
}

func (s *Service) SetAsset(ctx context.Context, asset *model.Asset) error {
	cacheKey := fmt.Sprintf("asset:%s", asset.Slug)

	data, err := json.Marshal(asset)
	if err != nil {
		return errors.NewInternalError("failed to marshal asset", err)
	}

	ttl := 30 * time.Minute
	if err := s.client.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return errors.NewTransientError("failed to set asset in cache", err)
	}

	return nil
}

func (s *Service) GetAssets(ctx context.Context) ([]model.Asset, error) {
	cacheKey := "assets:list"

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.NewTransientError("failed to get assets from cache", err)
	}

	var assets []model.Asset
	if err := json.Unmarshal(data, &assets); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal assets", err)
	}

	return assets, nil
}

func (s *Service) SetAssets(ctx context.Context, assets []model.Asset) error {
	cacheKey := "assets:list"

	data, err := json.Marshal(assets)
	if err != nil {
		return errors.NewInternalError("failed to marshal assets", err)
	}

	ttl := 30 * time.Minute
	if err := s.client.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
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
