package security

import (
	"context"
	"sync"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type requestBucket struct {
	timestamps []time.Time
}

// InMemoryRateLimiter is a simple sliding-window limiter intended for development / low-traffic deployments.
// It keeps at most `limit` timestamps per key and cleans buckets periodically so memory usage is bounded.

type InMemoryRateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*requestBucket
	limit    int
	window   time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

func NewInMemoryRateLimiter(limit int, window time.Duration) *InMemoryRateLimiter {
	rl := &InMemoryRateLimiter{
		buckets:  make(map[string]*requestBucket),
		limit:    limit,
		window:   window,
		cleanup:  5 * time.Minute,
		stopChan: make(chan struct{}),
	}

	go rl.cleanupLoop()
	return rl
}

func (rl *InMemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok := rl.buckets[key]
	if !ok {
		bucket = &requestBucket{}
		rl.buckets[key] = bucket
	}

	// filter timestamps within window
	valid := bucket.timestamps[:0]
	for _, ts := range bucket.timestamps {
		if ts.After(windowStart) {
			valid = append(valid, ts)
		}
	}
	bucket.timestamps = valid

	if len(bucket.timestamps) >= rl.limit {
		return false, nil
	}

	bucket.timestamps = append(bucket.timestamps, now)
	return true, nil
}

func (rl *InMemoryRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopChan:
			return
		case <-ticker.C:
			rl.purge()
		}
	}
}

func (rl *InMemoryRateLimiter) purge() {
	cutoff := time.Now().Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key, bucket := range rl.buckets {
		valid := bucket.timestamps[:0]
		for _, ts := range bucket.timestamps {
			if ts.After(cutoff) {
				valid = append(valid, ts)
			}
		}
		if len(valid) == 0 {
			delete(rl.buckets, key)
		} else {
			bucket.timestamps = valid
		}
	}
}

// Close stops the background cleanup goroutine (mainly for tests).
func (rl *InMemoryRateLimiter) Close() {
	close(rl.stopChan)
}
