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

	var bucketData bucket.Bucket
	if err := json.Unmarshal(data, &bucketData); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal bucket", err)
	}

	return &bucketData, nil
}

func (s *Service) SetBucket(ctx context.Context, bucket *bucket.Bucket) error {
	cacheKey := fmt.Sprintf("bucket:%s", bucket.Key())

	data, err := json.Marshal(bucket)
	if err != nil {
		return errors.NewInternalError("failed to marshal bucket", err)
	}

	if err := s.client.client.Set(ctx, cacheKey, data, s.ttl.Bucket).Err(); err != nil {
		return errors.NewTransientError("failed to set bucket in cache", err)
	}

	if err := s.client.client.Set(ctx, "bucket:keys", bucket.Key(), s.ttl.Bucket).Err(); err != nil {
		s.client.logger.WithError(err).Warn("Failed to cache bucket key for invalidation")
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

	var buckets []*bucket.Bucket
	if err := json.Unmarshal(data, &buckets); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal buckets", err)
	}

	var filtered []*bucket.Bucket
	for i, b := range buckets {
		if b == nil {
			s.client.logger.Warn(fmt.Sprintf("Nil bucket found in cache at index %d!", i))
			continue
		}
		filtered = append(filtered, b)
	}

	start := 0
	if nextKey != nil {
		for i, b := range filtered {
			if b.ID().Value() == *nextKey {
				start = i + 1
				break
			}
		}
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], nil
}

func (s *Service) SetBuckets(ctx context.Context, buckets []*bucket.Bucket) error {
	cacheKey := "buckets:list"

	// Filter out nil buckets and log if any found
	var filtered []*bucket.Bucket
	for i, b := range buckets {
		if b == nil {
			s.client.logger.Warn(fmt.Sprintf("Nil bucket found before caching at index %d!", i))
			continue
		}
		filtered = append(filtered, b)
	}

	data, err := json.Marshal(filtered)
	if err != nil {
		return errors.NewInternalError("failed to marshal buckets", err)
	}

	if err := s.client.client.Set(ctx, cacheKey, data, s.ttl.BucketsList).Err(); err != nil {
		return errors.NewTransientError("failed to set buckets in cache", err)
	}

	return nil
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

	var assetData asset.Asset
	if err := json.Unmarshal(data, &assetData); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal asset", err)
	}

	return &assetData, nil
}

func (s *Service) SetAsset(ctx context.Context, asset *asset.Asset) error {
	cacheKey := fmt.Sprintf("asset:%s", asset.Slug())

	data, err := json.Marshal(asset)
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

	var assets []*asset.Asset
	if err := json.Unmarshal(data, &assets); err != nil {
		return nil, errors.NewInternalError("failed to unmarshal assets", err)
	}

	return assets, nil
}

func (s *Service) SetAssets(ctx context.Context, assets []*asset.Asset) error {
	cacheKey := "assets:list"

	data, err := json.Marshal(assets)
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
