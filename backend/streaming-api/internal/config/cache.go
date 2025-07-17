package config

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/cache"
)

type CacheConfig struct {
	Service *cache.Service
	Client  *cache.Client
}

func NewCacheConfig(configManager *config.Manager, log *logger.Logger) (*CacheConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	host := dynamicCfg.GetStringFromComponent("redis", "host")
	port := dynamicCfg.GetIntFromComponent("redis", "port")
	db := dynamicCfg.GetIntFromComponent("redis", "db")
	password := dynamicCfg.GetStringFromComponent("redis", "password")

	log.Info("Redis config", "host", host, "port", port, "db", db, "password_set", password != "")

	redisClient, err := cache.NewRedisClientWithConfig(host, port, db, password)
	if err != nil {
		log.WithError(err).Error("Failed to connect to Redis")
		return nil, err
	}

	log.Info("Redis connection established", "host", host, "port", port, "db", db)

	cacheService := cache.NewService(redisClient)

	return &CacheConfig{
		Service: cacheService,
		Client:  redisClient,
	}, nil
}

func (cc *CacheConfig) Close() {
	if cc.Client != nil {
		cc.Client.Close()
	}
}
