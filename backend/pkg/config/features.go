package config

import (
	"sync"
	"time"
)

type FeatureFlag struct {
	Name        string
	Description string
	Enabled     bool
	Percentage  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FeatureFlagsManager struct {
	flags map[string]*FeatureFlag
	mu    sync.RWMutex
}

func NewFeatureFlagsManager() *FeatureFlagsManager {
	return &FeatureFlagsManager{
		flags: make(map[string]*FeatureFlag),
	}
}

func (ffm *FeatureFlagsManager) Register(name, description string, enabled bool) {
	ffm.mu.Lock()
	defer ffm.mu.Unlock()

	now := time.Now()
	ffm.flags[name] = &FeatureFlag{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		Percentage:  100,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (ffm *FeatureFlagsManager) IsEnabled(name string) bool {
	ffm.mu.RLock()
	defer ffm.mu.RUnlock()

	flag, exists := ffm.flags[name]
	if !exists {
		return false
	}

	return flag.Enabled
}

func (ffm *FeatureFlagsManager) IsEnabledForPercentage(name string, userID string) bool {
	ffm.mu.RLock()
	defer ffm.mu.RUnlock()

	flag, exists := ffm.flags[name]
	if !exists || !flag.Enabled {
		return false
	}

	if flag.Percentage >= 100 {
		return true
	}

	hash := hashString(userID)
	return hash%100 < flag.Percentage
}

func (ffm *FeatureFlagsManager) Enable(name string) {
	ffm.mu.Lock()
	defer ffm.mu.Unlock()

	if flag, exists := ffm.flags[name]; exists {
		flag.Enabled = true
		flag.UpdatedAt = time.Now()
	}
}

func (ffm *FeatureFlagsManager) Disable(name string) {
	ffm.mu.Lock()
	defer ffm.mu.Unlock()

	if flag, exists := ffm.flags[name]; exists {
		flag.Enabled = false
		flag.UpdatedAt = time.Now()
	}
}

func (ffm *FeatureFlagsManager) SetPercentage(name string, percentage int) {
	ffm.mu.Lock()
	defer ffm.mu.Unlock()

	if flag, exists := ffm.flags[name]; exists {
		if percentage < 0 {
			percentage = 0
		} else if percentage > 100 {
			percentage = 100
		}
		flag.Percentage = percentage
		flag.UpdatedAt = time.Now()
	}
}

func (ffm *FeatureFlagsManager) Get(name string) *FeatureFlag {
	ffm.mu.RLock()
	defer ffm.mu.RUnlock()

	return ffm.flags[name]
}

func (ffm *FeatureFlagsManager) List() []*FeatureFlag {
	ffm.mu.RLock()
	defer ffm.mu.RUnlock()

	flags := make([]*FeatureFlag, 0, len(ffm.flags))
	for _, flag := range ffm.flags {
		flags = append(flags, flag)
	}
	return flags
}

func (ffm *FeatureFlagsManager) Watch(configManager *Manager) {
	configManager.OnConfigChange(func(config *BaseConfig) {
		ffm.mu.Lock()
		defer ffm.mu.Unlock()

		if config.Features.EnableCircuitBreaker {
			ffm.enableIfNotExists("circuit_breaker", "Enable circuit breaker pattern")
		} else {
			ffm.disableIfExists("circuit_breaker")
		}

		if config.Features.EnableRetry {
			ffm.enableIfNotExists("retry", "Enable retry mechanism")
		} else {
			ffm.disableIfExists("retry")
		}

		if config.Features.EnableCaching {
			ffm.enableIfNotExists("caching", "Enable caching layer")
		} else {
			ffm.disableIfExists("caching")
		}

		if config.Features.EnableMetrics {
			ffm.enableIfNotExists("metrics", "Enable metrics collection")
		} else {
			ffm.disableIfExists("metrics")
		}

		if config.Features.EnableTracing {
			ffm.enableIfNotExists("tracing", "Enable distributed tracing")
		} else {
			ffm.disableIfExists("tracing")
		}
	})
}

func (ffm *FeatureFlagsManager) enableIfNotExists(name, description string) {
	if _, exists := ffm.flags[name]; !exists {
		now := time.Now()
		ffm.flags[name] = &FeatureFlag{
			Name:        name,
			Description: description,
			Enabled:     true,
			Percentage:  100,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	} else {
		ffm.flags[name].Enabled = true
		ffm.flags[name].UpdatedAt = time.Now()
	}
}

func (ffm *FeatureFlagsManager) disableIfExists(name string) {
	if flag, exists := ffm.flags[name]; exists {
		flag.Enabled = false
		flag.UpdatedAt = time.Now()
	}
}

func hashString(s string) int {
	hash := 0
	for _, char := range s {
		hash = ((hash << 5) - hash) + int(char)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}
