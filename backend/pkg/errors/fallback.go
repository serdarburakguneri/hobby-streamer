package errors

import (
	"context"
	"time"
)

type FallbackFunc func(ctx context.Context) error
type FallbackResult struct {
	Success bool
	Error   error
	Used    string
}

type FallbackChain struct {
	primary   FallbackFunc
	fallbacks []FallbackFunc
	names     []string
	timeout   time.Duration
}

func NewFallbackChain(primary FallbackFunc, primaryName string) *FallbackChain {
	return &FallbackChain{
		primary:   primary,
		fallbacks: make([]FallbackFunc, 0),
		names:     []string{primaryName},
		timeout:   5 * time.Second,
	}
}

func (fc *FallbackChain) AddFallback(fn FallbackFunc, name string) *FallbackChain {
	fc.fallbacks = append(fc.fallbacks, fn)
	fc.names = append(fc.names, name)
	return fc
}

func (fc *FallbackChain) SetTimeout(timeout time.Duration) *FallbackChain {
	fc.timeout = timeout
	return fc
}

func (fc *FallbackChain) Execute(ctx context.Context) *FallbackResult {
	ctx, cancel := context.WithTimeout(ctx, fc.timeout)
	defer cancel()

	if fc.primary != nil {
		err := fc.primary(ctx)
		if err == nil {
			return &FallbackResult{
				Success: true,
				Error:   nil,
				Used:    fc.names[0],
			}
		}
	}

	for i, fallback := range fc.fallbacks {
		select {
		case <-ctx.Done():
			return &FallbackResult{
				Success: false,
				Error:   ctx.Err(),
				Used:    "timeout",
			}
		default:
		}

		err := fallback(ctx)
		if err == nil {
			return &FallbackResult{
				Success: true,
				Error:   nil,
				Used:    fc.names[i+1],
			}
		}
	}

	return &FallbackResult{
		Success: false,
		Error:   NewInternalError("all fallback options failed", nil),
		Used:    "none",
	}
}

type CacheFallback struct {
	cache    interface{}
	fallback interface{}
}

func NewCacheFallback(cache, fallback interface{}) *CacheFallback {
	return &CacheFallback{
		cache:    cache,
		fallback: fallback,
	}
}

func (cf *CacheFallback) Get() interface{} {
	if cf.cache != nil {
		return cf.cache
	}
	return cf.fallback
}

func (cf *CacheFallback) SetCache(value interface{}) {
	cf.cache = value
}

func (cf *CacheFallback) SetFallback(value interface{}) {
	cf.fallback = value
}

type DegradationLevel int

const (
	DegradationNone DegradationLevel = iota
	DegradationPartial
	DegradationFull
)

type DegradationManager struct {
	level    DegradationLevel
	handlers map[DegradationLevel]func()
}

func NewDegradationManager() *DegradationManager {
	return &DegradationManager{
		level:    DegradationNone,
		handlers: make(map[DegradationLevel]func()),
	}
}

func (dm *DegradationManager) SetLevel(level DegradationLevel) {
	if dm.level != level {
		dm.level = level
		if handler, exists := dm.handlers[level]; exists {
			handler()
		}
	}
}

func (dm *DegradationManager) GetLevel() DegradationLevel {
	return dm.level
}

func (dm *DegradationManager) OnLevelChange(level DegradationLevel, handler func()) {
	dm.handlers[level] = handler
}

func (dm *DegradationManager) IsDegraded() bool {
	return dm.level > DegradationNone
}

func (dm *DegradationManager) IsFullyDegraded() bool {
	return dm.level == DegradationFull
}
