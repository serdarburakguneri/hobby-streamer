package security

import (
	"net/http"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func RateLimitMiddleware(limit int, windowDuration time.Duration) func(http.Handler) http.Handler {
	limiter := NewInMemoryRateLimiter(limit, windowDuration)
	return RateLimitMiddlewareWithLimiter(limiter)
}

func RateLimitMiddlewareWithLimiter(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getClientIP(r)

			allowed, err := limiter.Allow(r.Context(), key)
			if err != nil {
				logger.Get().WithError(err).Error("rate limiter error")
				next.ServeHTTP(w, r)
				return
			}
			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error": "Rate limit exceeded"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
