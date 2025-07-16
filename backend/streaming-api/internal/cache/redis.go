package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
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

func NewRedisClient() (*Client, error) {
	redisURL := getRedisURL()

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
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
		return nil, fmt.Errorf("failed to get bucket from cache: %w", err)
	}

	var bucket model.Bucket
	if err := json.Unmarshal(data, &bucket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bucket: %w", err)
	}

	return &bucket, nil
}

func (s *Service) SetBucket(ctx context.Context, bucket *model.Bucket) error {
	cacheKey := fmt.Sprintf("bucket:%s", bucket.Key)

	data, err := json.Marshal(bucket)
	if err != nil {
		return fmt.Errorf("failed to marshal bucket: %w", err)
	}

	ttl := 30 * time.Minute
	if err := s.client.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set bucket in cache: %w", err)
	}

	return nil
}

func (s *Service) GetBuckets(ctx context.Context) ([]model.Bucket, error) {
	data, err := s.client.client.Get(ctx, "buckets:list").Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get buckets from cache: %w", err)
	}

	var buckets []model.Bucket
	if err := json.Unmarshal(data, &buckets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal buckets: %w", err)
	}

	return buckets, nil
}

func (s *Service) SetBuckets(ctx context.Context, buckets []model.Bucket) error {
	data, err := json.Marshal(buckets)
	if err != nil {
		return fmt.Errorf("failed to marshal buckets: %w", err)
	}

	ttl := 15 * time.Minute
	if err := s.client.client.Set(ctx, "buckets:list", data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set buckets in cache: %w", err)
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
		return nil, fmt.Errorf("failed to get asset from cache: %w", err)
	}

	var asset model.Asset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}

	return &asset, nil
}

func (s *Service) SetAsset(ctx context.Context, asset *model.Asset) error {
	cacheKey := fmt.Sprintf("asset:%s", asset.Slug)

	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}

	ttl := 30 * time.Minute
	if err := s.client.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set asset in cache: %w", err)
	}

	return nil
}

func (s *Service) GetAssets(ctx context.Context) ([]model.Asset, error) {
	data, err := s.client.client.Get(ctx, "assets:list").Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get assets from cache: %w", err)
	}

	var assets []model.Asset
	if err := json.Unmarshal(data, &assets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal assets: %w", err)
	}

	return assets, nil
}

func (s *Service) SetAssets(ctx context.Context, assets []model.Asset) error {
	data, err := json.Marshal(assets)
	if err != nil {
		return fmt.Errorf("failed to marshal assets: %w", err)
	}

	ttl := 15 * time.Minute
	if err := s.client.client.Set(ctx, "assets:list", data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set assets in cache: %w", err)
	}

	return nil
}

func (s *Service) InvalidateBucket(ctx context.Context, key string) error {
	cacheKey := fmt.Sprintf("bucket:%s", key)
	if err := s.client.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("failed to invalidate bucket cache: %w", err)
	}

	if err := s.client.client.Del(ctx, "buckets:list").Err(); err != nil {
		return fmt.Errorf("failed to invalidate buckets list cache: %w", err)
	}

	return nil
}

func (s *Service) InvalidateAsset(ctx context.Context, slug string) error {
	cacheKey := fmt.Sprintf("asset:%s", slug)
	if err := s.client.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("failed to invalidate asset cache: %w", err)
	}

	if err := s.client.client.Del(ctx, "assets:list").Err(); err != nil {
		return fmt.Errorf("failed to invalidate assets list cache: %w", err)
	}

	return nil
}

func getRedisURL() string {
	if url := getEnv("REDIS_URL", ""); url != "" {
		return url
	}

	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")

	if password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s", password, host, port)
	}

	return fmt.Sprintf("redis://%s:%s", host, port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
