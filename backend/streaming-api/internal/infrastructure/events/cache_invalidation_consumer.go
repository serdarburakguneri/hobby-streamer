package events

import (
	"context"
	"encoding/json"

	pkgevents "github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	streamcache "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/cache"
)

type CacheInvalidator struct {
	consumer *pkgevents.Consumer
	cache    streamcache.CacheService
	logger   *logger.Logger
}

func NewCacheInvalidator(cacheSvc streamcache.CacheService) *CacheInvalidator {
	return &CacheInvalidator{cache: cacheSvc, logger: logger.WithService("streaming-cache-invalidator")}
}

func (c *CacheInvalidator) Start(ctx context.Context, bootstrapServers string) error {
	cfg := pkgevents.DefaultConsumerConfig()
	cfg.BootstrapServers = []string{bootstrapServers}
	cfg.GroupID = "streaming-api-cache-group"
	cfg.Topics = []string{pkgevents.AssetEventsTopic, pkgevents.BucketEventsTopic}

	cons, err := pkgevents.NewConsumer(ctx, cfg)
	if err != nil {
		return err
	}
	c.consumer = cons

	cons.Subscribe(pkgevents.AssetEventsTopic, c.handleAssetEvent)
	cons.Subscribe(pkgevents.BucketEventsTopic, c.handleBucketEvent)

	go func() {
		if err := cons.Start(ctx); err != nil {
			c.logger.WithError(err).Error("Kafka consumer error")
		}
	}()
	c.logger.Info("Started streaming-api cache invalidator consumer")
	return nil
}

func (c *CacheInvalidator) Stop() error {
	if c.consumer != nil {
		return c.consumer.Stop()
	}
	return nil
}

func (c *CacheInvalidator) handleAssetEvent(ctx context.Context, ev *pkgevents.Event) error {
	c.logger.Info("Asset event received", "type", ev.Type)

	switch ev.Type {
	case pkgevents.VideoStatusUpdatedEventType,
		pkgevents.AssetUpdatedEventType,
		pkgevents.AssetCreatedEventType,
		pkgevents.AssetDeletedEventType,
		pkgevents.VideoAddedEventType,
		pkgevents.VideoRemovedEventType:
		c.cache.InvalidateAssetsListCache(ctx)

		if slug := extractString(ev.Data, "slug"); slug != "" {
			c.cache.InvalidateAssetCache(ctx, slug)
		}
	}
	return nil
}

func (c *CacheInvalidator) handleBucketEvent(ctx context.Context, ev *pkgevents.Event) error {
	c.logger.Info("Bucket event received", "type", ev.Type)
	switch ev.Type {
	case pkgevents.BucketUpdatedEventType,
		pkgevents.BucketCreatedEventType,
		pkgevents.BucketDeletedEventType,
		pkgevents.BucketAssetAddedEventType,
		pkgevents.BucketAssetRemovedEventType:
		c.cache.InvalidateBucketsListCache(ctx)
		if key := extractString(ev.Data, "key"); key != "" {
			c.cache.InvalidateBucketCache(ctx, key)
		}
	}
	return nil
}

func extractString(data interface{}, field string) string {
	switch v := data.(type) {
	case map[string]interface{}:
		if s, ok := v[field].(string); ok {
			return s
		}
		b, _ := json.Marshal(v)
		var m map[string]interface{}
		_ = json.Unmarshal(b, &m)
		if s, ok := m[field].(string); ok {
			return s
		}
	case []byte:
		var m map[string]interface{}
		if err := json.Unmarshal(v, &m); err == nil {
			if s, ok := m[field].(string); ok {
				return s
			}
		}
	}
	return ""
}
