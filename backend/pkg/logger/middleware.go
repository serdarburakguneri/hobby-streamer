package logger

import (
	"context"
	"net/http"
	"time"
)

func RequestLoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			trackingID := r.Header.Get("X-Tracking-ID")
			if trackingID == "" {
				trackingID = GenerateTrackingID()
			}

			w.Header().Set("X-Tracking-ID", trackingID)

			ctx := context.WithValue(r.Context(), "tracking_id", trackingID)
			r = r.WithContext(ctx)

			trackedLogger := logger.WithTrackingID(trackingID)

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			trackedLogger.LogRequest(r, wrapped.statusCode, duration)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
