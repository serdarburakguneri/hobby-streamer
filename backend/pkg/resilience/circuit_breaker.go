package resilience

import (
	"context"
	"sync"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	name          string
	state         CircuitState
	failureCount  int64
	lastFailure   time.Time
	threshold     int64
	timeout       time.Duration
	mutex         sync.RWMutex
	onStateChange func(name string, from, to CircuitState)
}

type CircuitBreakerConfig struct {
	Name          string
	Threshold     int64
	Timeout       time.Duration
	OnStateChange func(name string, from, to CircuitState)
}

func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.Threshold <= 0 {
		config.Threshold = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	return &CircuitBreaker{
		name:          config.Name,
		state:         CircuitClosed,
		threshold:     config.Threshold,
		timeout:       config.Timeout,
		onStateChange: config.OnStateChange,
	}
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if !cb.Ready() {
		return pkgerrors.NewCircuitBreakerError("circuit breaker is open", nil)
	}

	err := fn()
	cb.RecordResult(err)
	return err
}

func (cb *CircuitBreaker) Ready() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailure) >= cb.timeout {
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.transitionTo(CircuitHalfOpen)
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) RecordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err == nil {
		cb.onSuccess()
	} else {
		cb.onFailure()
	}
}

func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case CircuitClosed:
		cb.failureCount = 0
	case CircuitHalfOpen:
		cb.transitionTo(CircuitClosed)
		cb.failureCount = 0
	}
}

func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailure = time.Now()

	switch cb.state {
	case CircuitClosed:
		if cb.failureCount >= cb.threshold {
			cb.transitionTo(CircuitOpen)
		}
	case CircuitHalfOpen:
		cb.transitionTo(CircuitOpen)
	}
}

func (cb *CircuitBreaker) transitionTo(newState CircuitState) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, oldState, newState)
	}
}

func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.transitionTo(CircuitClosed)
	cb.failureCount = 0
}

func (cb *CircuitBreaker) ForceOpen() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.transitionTo(CircuitOpen)
}

func (cb *CircuitBreaker) ForceClose() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.transitionTo(CircuitClosed)
	cb.failureCount = 0
}

type CircuitBreakerRegistry struct {
	breakers map[string]*CircuitBreaker
	mutex    sync.RWMutex
}

func NewCircuitBreakerRegistry() *CircuitBreakerRegistry {
	return &CircuitBreakerRegistry{
		breakers: make(map[string]*CircuitBreaker),
	}
}

func (r *CircuitBreakerRegistry) GetOrCreate(name string, config CircuitBreakerConfig) *CircuitBreaker {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if breaker, exists := r.breakers[name]; exists {
		return breaker
	}

	config.Name = name
	breaker := NewCircuitBreaker(config)
	r.breakers[name] = breaker
	return breaker
}

func (r *CircuitBreakerRegistry) Get(name string) *CircuitBreaker {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.breakers[name]
}

func (r *CircuitBreakerRegistry) Reset(name string) {
	if breaker := r.Get(name); breaker != nil {
		breaker.Reset()
	}
}

func (r *CircuitBreakerRegistry) ResetAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, breaker := range r.breakers {
		breaker.Reset()
	}
}
