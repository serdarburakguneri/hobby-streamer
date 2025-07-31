package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
)

type Client struct {
	client *redis.Client
	logger *logger.Logger
}

type Service struct {
	client     *Client
	marshaller *DomainMarshaller
	ttl        TTLConfig
}

type TTLConfig struct {
	Bucket      time.Duration
	BucketsList time.Duration
	Asset       time.Duration
	AssetsList  time.Duration
}

func NewRedisClientWithConfig(host string, port int, db int, password string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "redis_ping",
			"host":      host,
			"port":      port,
			"db":        db,
		})
	}

	return &Client{
		client: client,
		logger: logger.Get().WithService("redis-client"),
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func NewService(client *Client, ttl TTLConfig) *Service {
	return &Service{
		client:     client,
		marshaller: NewDomainMarshaller(),
		ttl:        ttl,
	}
}

func (s *Service) GetBucket(ctx context.Context, key string) (*bucketentity.Bucket, error) {
	cacheKey := s.marshaller.GenerateBucketKey(key)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "redis_get_bucket",
			"key":       key,
		})
	}

	return s.marshaller.UnmarshalBucket(data)
}

func (s *Service) SetBucket(ctx context.Context, bucket *bucketentity.Bucket) error {
	if bucket == nil {
		return errors.NewValidationError("bucket cannot be nil", nil)
	}

	data, err := s.marshaller.MarshalBucket(bucket)
	if err != nil {
		return err
	}

	cacheKey := s.marshaller.GenerateBucketKey(bucket.Key().Value())

	err = s.client.client.Set(ctx, cacheKey, data, s.ttl.Bucket).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_set_bucket",
			"key":       bucket.Key().Value(),
		})
	}

	return nil
}

func (s *Service) GetBuckets(ctx context.Context, limit int, nextKey *string) ([]*bucketentity.Bucket, error) {
	cacheKey := s.marshaller.GenerateBucketsListKey(limit, nextKey)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "redis_get_buckets",
			"limit":     limit,
		})
	}

	return s.marshaller.UnmarshalBuckets(data)
}

func (s *Service) SetBuckets(ctx context.Context, buckets []*bucketentity.Bucket, limit int, nextKey *string) error {
	if buckets == nil {
		return errors.NewValidationError("buckets cannot be nil", nil)
	}

	data, err := s.marshaller.MarshalBuckets(buckets)
	if err != nil {
		return err
	}

	cacheKey := s.marshaller.GenerateBucketsListKey(limit, nextKey)

	err = s.client.client.Set(ctx, cacheKey, data, s.ttl.BucketsList).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_set_buckets",
			"count":     len(buckets),
		})
	}

	return nil
}

func (s *Service) GetAsset(ctx context.Context, slug string) (*assetentity.Asset, error) {
	cacheKey := s.marshaller.GenerateAssetKey(slug)

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "redis_get_asset",
			"slug":      slug,
		})
	}

	return s.marshaller.UnmarshalAsset(data)
}

func (s *Service) SetAsset(ctx context.Context, asset *assetentity.Asset) error {
	if asset == nil {
		return errors.NewValidationError("asset cannot be nil", nil)
	}

	data, err := s.marshaller.MarshalAsset(asset)
	if err != nil {
		return err
	}

	cacheKey := s.marshaller.GenerateAssetKey(asset.Slug().Value())

	err = s.client.client.Set(ctx, cacheKey, data, s.ttl.Asset).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_set_asset",
			"slug":      asset.Slug().Value(),
		})
	}

	return nil
}

func (s *Service) GetAssets(ctx context.Context) ([]*assetentity.Asset, error) {
	cacheKey := s.marshaller.GenerateAssetsListKey()

	data, err := s.client.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "redis_get_assets",
		})
	}

	return s.marshaller.UnmarshalAssets(data)
}

func (s *Service) SetAssets(ctx context.Context, assets []*assetentity.Asset) error {
	if assets == nil {
		return errors.NewValidationError("assets cannot be nil", nil)
	}

	data, err := s.marshaller.MarshalAssets(assets)
	if err != nil {
		return err
	}

	cacheKey := s.marshaller.GenerateAssetsListKey()

	err = s.client.client.Set(ctx, cacheKey, data, s.ttl.AssetsList).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_set_assets",
			"count":     len(assets),
		})
	}

	return nil
}

func (s *Service) InvalidateBucketCache(ctx context.Context, key string) error {
	cacheKey := s.marshaller.GenerateBucketKey(key)

	err := s.client.client.Del(ctx, cacheKey).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_invalidate_bucket",
			"key":       key,
		})
	}

	return nil
}

func (s *Service) InvalidateBucketsListCache(ctx context.Context) error {
	pattern := "buckets:list:*"

	keys, err := s.client.client.Keys(ctx, pattern).Result()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_invalidate_buckets_list",
		})
	}

	if len(keys) > 0 {
		err = s.client.client.Del(ctx, keys...).Err()
		if err != nil {
			return errors.WithContext(err, map[string]interface{}{
				"operation": "redis_delete_buckets_list_keys",
				"count":     len(keys),
			})
		}
	}

	return nil
}

func (s *Service) InvalidateAssetCache(ctx context.Context, slug string) error {
	cacheKey := s.marshaller.GenerateAssetKey(slug)

	err := s.client.client.Del(ctx, cacheKey).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_invalidate_asset",
			"slug":      slug,
		})
	}

	return nil
}

func (s *Service) InvalidateAssetsListCache(ctx context.Context) error {
	cacheKey := s.marshaller.GenerateAssetsListKey()

	err := s.client.client.Del(ctx, cacheKey).Err()
	if err != nil {
		return errors.WithContext(err, map[string]interface{}{
			"operation": "redis_invalidate_assets_list",
		})
	}

	return nil
}
