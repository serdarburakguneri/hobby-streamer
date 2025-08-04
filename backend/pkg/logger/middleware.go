package logger

import (
	"compress/gzip"
	"context"
	"net/http"
	"strings"
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

			ctx := context.WithValue(r.Context(), trackingIDContextKey, trackingID)
			r = r.WithContext(ctx)

			trackedLogger := logger.WithTrackingID(trackingID)

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			trackedLogger.LogRequest(r, wrapped.statusCode, duration)
		})
	}
}

func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		next.ServeHTTP(&gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter:     gzipWriter,
		}, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (gzw *gzipResponseWriter) Write(data []byte) (int, error) {
	return gzw.gzipWriter.Write(data)
}

func (gzw *gzipResponseWriter) WriteString(s string) (int, error) {
	return gzw.gzipWriter.Write([]byte(s))
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
