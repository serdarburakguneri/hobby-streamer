package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type RedisRateLimiter struct {
	client *redis.Client
	logger *logger.Logger
}

func NewRedisRateLimiter(addr, password string, db int) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisRateLimiter{
		client: client,
		logger: logger.WithService("redis-rate-limiter"),
	}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-window)

	script := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local limit = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		local count = redis.call('ZCARD', key)
		
		if count >= limit then
			return 0
		end
		
		redis.call('ZADD', key, now, now .. ':' .. math.random())
		redis.call('EXPIRE', key, ARGV[4])
		return 1
	`

	result, err := rl.client.Eval(ctx, script, []string{key}, windowStart.Unix(), limit, now.Unix(), int(window.Seconds())).Result()
	if err != nil {
		rl.logger.WithError(err).Error("Failed to execute rate limit script", "key", key)
		return false, err
	}

	allowed := result.(int64) == 1
	if !allowed {
		rl.logger.Debug("Rate limit exceeded", "key", key, "limit", limit, "window", window)
	}

	return allowed, nil
}

func (rl *RedisRateLimiter) Close() error {
	return rl.client.Close()
}

func (rl *RedisRateLimiter) Ping(ctx context.Context) error {
	return rl.client.Ping(ctx).Err()
}
