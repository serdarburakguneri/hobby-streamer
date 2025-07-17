package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger
}

func New(level slog.Level, format string) *Logger {
	var handler slog.Handler

	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

func GenerateTrackingID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (l *Logger) WithTrackingID(trackingID string) *Logger {
	return &Logger{
		Logger: l.Logger.With("tracking_id", trackingID),
	}
}

func (l *Logger) WithService(service string) *Logger {
	return &Logger{
		Logger: l.Logger.With("service", service),
	}
}

func (l *Logger) WithRequest(r *http.Request) *Logger {
	attrs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	}

	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}

	if trackingID := r.Header.Get("X-Tracking-ID"); trackingID != "" {
		attrs = append(attrs, "tracking_id", trackingID)
	}

	if user := r.Context().Value("user"); user != nil {
		// Try to extract user info from context
		if userMap, ok := user.(map[string]interface{}); ok {
			if id, exists := userMap["id"]; exists {
				attrs = append(attrs, "user_id", id)
			}
			if username, exists := userMap["username"]; exists {
				attrs = append(attrs, "username", username)
			}
		}
	}

	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	attrs := []any{}

	if requestID := ctx.Value("request_id"); requestID != nil {
		attrs = append(attrs, "request_id", requestID)
	}

	if trackingID := ctx.Value("tracking_id"); trackingID != nil {
		attrs = append(attrs, "tracking_id", trackingID)
	}

	if user := ctx.Value("user"); user != nil {
		if userMap, ok := user.(map[string]interface{}); ok {
			if id, exists := userMap["id"]; exists {
				attrs = append(attrs, "user_id", id)
			}
			if username, exists := userMap["username"]; exists {
				attrs = append(attrs, "username", username)
			}
		}
	}

	if len(attrs) == 0 {
		return l
	}

	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With("error", err.Error()),
	}
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}

	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *Logger) LogRequest(r *http.Request, statusCode int, duration time.Duration) {
	level := slog.LevelInfo
	if statusCode >= 400 {
		level = slog.LevelWarn
	}
	if statusCode >= 500 {
		level = slog.LevelError
	}

	logger := l.WithRequest(r)
	switch level {
	case slog.LevelError:
		logger.Error("HTTP request completed", "status_code", statusCode, "duration_ms", duration.Milliseconds())
	case slog.LevelWarn:
		logger.Warn("HTTP request completed", "status_code", statusCode, "duration_ms", duration.Milliseconds())
	default:
		logger.Info("HTTP request completed", "status_code", statusCode, "duration_ms", duration.Milliseconds())
	}
}

func (l *Logger) LogError(err error, msg string, args ...any) {
	allArgs := append([]any{"error", err.Error()}, args...)
	l.Logger.Error(msg, allArgs...)
}

var defaultLogger *Logger

func Init(level slog.Level, format string) {
	defaultLogger = New(level, format)
}

func Get() *Logger {
	if defaultLogger == nil {
		Init(slog.LevelInfo, "text")
	}
	return defaultLogger
}

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func WithService(service string) *Logger {
	return Get().WithService(service)
}

func WithTrackingID(trackingID string) *Logger {
	return Get().WithTrackingID(trackingID)
}

func WithContext(ctx context.Context) *Logger {
	return Get().WithContext(ctx)
}

func WithError(err error) *Logger {
	return Get().WithError(err)
}

func WithFields(fields map[string]any) *Logger {
	return Get().WithFields(fields)
}

func GetLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
