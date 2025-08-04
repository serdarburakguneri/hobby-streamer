package resilience

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	JitterFactor    float64
	RetryableErrors []pkgerrors.ErrorType
}

type RetryableFunc func(ctx context.Context) error

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		JitterFactor:  0.1,
		RetryableErrors: []pkgerrors.ErrorType{
			pkgerrors.ErrorTypeTransient,
			pkgerrors.ErrorTypeTimeout,
			pkgerrors.ErrorTypeExternal,
		},
	}
}

func (c *RetryConfig) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errorType := pkgerrors.GetErrorType(err)
	for _, retryableType := range c.RetryableErrors {
		if errorType == retryableType {
			return true
		}
	}

	return false
}

func (c *RetryConfig) calculateDelay(attempt int) time.Duration {
	delay := float64(c.InitialDelay) * math.Pow(c.BackoffFactor, float64(attempt-1))

	if delay > float64(c.MaxDelay) {
		delay = float64(c.MaxDelay)
	}

	if c.JitterFactor > 0 {
		var buf [8]byte
		if _, err := rand.Read(buf[:]); err == nil {
			randomFloat := float64(binary.BigEndian.Uint64(buf[:])) / float64(^uint64(0))
			jitter := delay * c.JitterFactor * (randomFloat*2 - 1)
			delay += jitter
			if delay < 0 {
				delay = 0
			}
		}
	}

	return time.Duration(delay)
}

func Retry(ctx context.Context, fn RetryableFunc, config *RetryConfig) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		if !config.IsRetryableError(err) {
			return err
		}

		if attempt == config.MaxAttempts {
			break
		}

		delay := config.calculateDelay(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return lastErr
}

func RetryWithBackoff(ctx context.Context, fn RetryableFunc, maxAttempts int) error {
	config := &RetryConfig{
		MaxAttempts:   maxAttempts,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		JitterFactor:  0.1,
		RetryableErrors: []pkgerrors.ErrorType{
			pkgerrors.ErrorTypeTransient,
			pkgerrors.ErrorTypeTimeout,
			pkgerrors.ErrorTypeExternal,
		},
	}
	return Retry(ctx, fn, config)
}

func RetryFast(ctx context.Context, fn RetryableFunc) error {
	config := &RetryConfig{
		MaxAttempts:   2,
		InitialDelay:  50 * time.Millisecond,
		MaxDelay:      200 * time.Millisecond,
		BackoffFactor: 2.0,
		JitterFactor:  0.1,
		RetryableErrors: []pkgerrors.ErrorType{
			pkgerrors.ErrorTypeTransient,
			pkgerrors.ErrorTypeTimeout,
		},
	}
	return Retry(ctx, fn, config)
}
