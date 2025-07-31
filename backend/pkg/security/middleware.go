package security

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type InMemoryRateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewInMemoryRateLimiter(limit int, window time.Duration) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *InMemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests

		if len(validRequests) >= rl.limit {
			return false, nil
		}
	}

	rl.requests[key] = append(rl.requests[key], now)
	return true, nil
}

func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewInMemoryRateLimiter(limit, window)
	return RateLimitMiddlewareWithLimiter(limiter)
}

func RateLimitMiddlewareWithLimiter(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getClientIP(r)

			allowed, err := limiter.Allow(r.Context(), key)
			if err != nil {
				logger.Get().WithError(err).Error("Rate limiter error")
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				if _, err := w.Write([]byte(`{"error": "Rate limit exceeded"}`)); err != nil {
					logger.Get().WithError(err).Error("Failed to write rate limit response")
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' ws: wss:;")

			next.ServeHTTP(w, r)
		})
	}
}

func CORSMiddleware(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) func(http.Handler) http.Handler {
	originSet := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originSet[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" && originSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func InputValidationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > 10*1024*1024 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				if _, err := w.Write([]byte(`{"error": "Request too large"}`)); err != nil {
					logger.Get().WithError(err).Error("Failed to write request too large response")
				}
				return
			}

			contentType := r.Header.Get("Content-Type")
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				if !strings.Contains(contentType, "application/json") &&
					!strings.Contains(contentType, "multipart/form-data") &&
					!strings.Contains(contentType, "application/x-www-form-urlencoded") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnsupportedMediaType)
					if _, err := w.Write([]byte(`{"error": "Unsupported content type"}`)); err != nil {
						logger.Get().WithError(err).Error("Failed to write unsupported content type response")
					}
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if comma := strings.Index(xff, ","); comma != -1 {
			return strings.TrimSpace(xff[:comma])
		}
		return strings.TrimSpace(xff)
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	return r.RemoteAddr
}

func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", getClientIP(r),
				"user_agent", r.UserAgent(),
				"content_length", r.ContentLength,
			)

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			log.Info("HTTP response",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration,
			)
		})
	}
}
