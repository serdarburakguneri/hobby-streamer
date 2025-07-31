package logger

import (
	"log/slog"
	"testing"
	"time"
)

func BenchmarkSyncLogger(b *testing.B) {
	logger := New(slog.LevelInfo, "text")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test message", "key", "value", "number", 42)
		}
	})
}

func BenchmarkAsyncLogger(b *testing.B) {
	logger := NewAsync(slog.LevelInfo, "text", 1000)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test message", "key", "value", "number", 42)
		}
	})
}

func BenchmarkSyncLoggerWithFields(b *testing.B) {
	logger := New(slog.LevelInfo, "text")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(map[string]any{
				"user_id": "123",
				"action":  "test",
				"time":    time.Now(),
			}).Info("test message", "key", "value")
		}
	})
}

func BenchmarkAsyncLoggerWithFields(b *testing.B) {
	logger := NewAsync(slog.LevelInfo, "text", 1000)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(map[string]any{
				"user_id": "123",
				"action":  "test",
				"time":    time.Now(),
			}).Info("test message", "key", "value")
		}
	})
}
