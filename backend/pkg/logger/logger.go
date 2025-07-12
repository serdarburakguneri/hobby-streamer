package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Logger wraps slog.Logger to provide additional context and convenience methods
type Logger struct {
	*slog.Logger
}

// New creates a new logger with the specified level and format
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

// WithService adds service name to the logger context
func (l *Logger) WithService(service string) *Logger {
	return &Logger{
		Logger: l.Logger.With("service", service),
	}
}

// WithRequest adds request context to the logger
func (l *Logger) WithRequest(r *http.Request) *Logger {
	attrs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	}

	// Add request ID if present
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}

	// Add user context if present
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

// WithContext adds context values to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	attrs := []any{}

	// Add request ID from context if present
	if requestID := ctx.Value("request_id"); requestID != nil {
		attrs = append(attrs, "request_id", requestID)
	}

	// Add user from context if present
	if user := ctx.Value("user"); user != nil {
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

	if len(attrs) == 0 {
		return l
	}

	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

// WithError adds error context to the logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With("error", err.Error()),
	}
}

// WithFields adds custom fields to the logger
func (l *Logger) WithFields(fields map[string]any) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}

	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

// Convenience methods for common logging patterns
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

// LogRequest logs HTTP request details
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

// LogError logs errors with context
func (l *Logger) LogError(err error, msg string, args ...any) {
	allArgs := append([]any{"error", err.Error()}, args...)
	l.Logger.Error(msg, allArgs...)
}

// Global logger instance
var defaultLogger *Logger

// Init initializes the global logger
func Init(level slog.Level, format string) {
	defaultLogger = New(level, format)
}

// Get returns the global logger instance
func Get() *Logger {
	if defaultLogger == nil {
		// Initialize with default settings if not already initialized
		Init(slog.LevelInfo, "text")
	}
	return defaultLogger
}

// Helper functions for global logger
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

func WithContext(ctx context.Context) *Logger {
	return Get().WithContext(ctx)
}

func WithError(err error) *Logger {
	return Get().WithError(err)
}

func WithFields(fields map[string]any) *Logger {
	return Get().WithFields(fields)
}
