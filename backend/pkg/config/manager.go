package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Manager struct {
	config    *BaseConfig
	viper     *viper.Viper
	mu        sync.RWMutex
	watcher   *fsnotify.Watcher
	callbacks []func(*BaseConfig)
}

func NewManager(serviceName string) (*Manager, error) {
	v := viper.New()

	env := getEnvironment()

	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("environment", env)
	v.SetDefault("service", serviceName)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("features.enable_circuit_breaker", true)
	v.SetDefault("features.enable_retry", true)
	v.SetDefault("features.enable_caching", true)
	v.SetDefault("features.enable_metrics", false)
	v.SetDefault("features.enable_tracing", false)
	v.SetDefault("circuit_breaker.threshold", 5)
	v.SetDefault("circuit_breaker.timeout", "30s")
	v.SetDefault("retry.max_attempts", 3)
	v.SetDefault("retry.base_delay", "100ms")
	v.SetDefault("retry.max_delay", "5s")
	v.SetDefault("cache.ttl", "30m")
	v.SetDefault("security.rate_limit.requests", 100)
	v.SetDefault("security.rate_limit.window", "1m")
	v.SetDefault("security.max_request_size", 10485760)
	v.SetDefault("security.cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:8081"})
	v.SetDefault("security.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("security.cors.allowed_headers", []string{"Content-Type", "Authorization", "X-Requested-With"})

	manager := &Manager{
		viper:     v,
		callbacks: make([]func(*BaseConfig), 0),
	}

	if err := manager.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return manager, nil
}

func (m *Manager) loadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config BaseConfig
	if err := m.viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := m.validateConfig(&config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	m.config = &config
	return nil
}

func (m *Manager) validateConfig(config *BaseConfig) error {
	if config.Service == "" {
		return fmt.Errorf("service name is required")
	}

	if config.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	if config.Log.Level == "" {
		return fmt.Errorf("log level is required")
	}

	if config.Log.Format == "" {
		return fmt.Errorf("log format is required")
	}

	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	return nil
}

func (m *Manager) GetConfig() *BaseConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

func (m *Manager) GetDynamicConfig() *DynamicConfig {
	return NewDynamicConfig(m.GetConfig())
}

func (m *Manager) WatchConfig() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	m.watcher = watcher

	configFile := m.viper.ConfigFileUsed()
	if configFile != "" {
		configDir := filepath.Dir(configFile)
		if err := watcher.Add(configDir); err != nil {
			return fmt.Errorf("failed to watch config directory: %w", err)
		}
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if strings.HasSuffix(event.Name, ".yaml") || strings.HasSuffix(event.Name, ".yml") {
						time.Sleep(100 * time.Millisecond)
						if err := m.loadConfig(); err != nil {
							fmt.Printf("Failed to reload config: %v\n", err)
							continue
						}
						m.notifyCallbacks()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Config watcher error: %v\n", err)
			}
		}
	}()

	return nil
}

func (m *Manager) OnConfigChange(callback func(*BaseConfig)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks = append(m.callbacks, callback)
}

func (m *Manager) notifyCallbacks() {
	m.mu.RLock()
	callbacks := make([]func(*BaseConfig), len(m.callbacks))
	copy(callbacks, m.callbacks)
	config := m.config
	m.mu.RUnlock()

	for _, callback := range callbacks {
		callback(config)
	}
}

func (m *Manager) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}

func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "development"
	}
	return strings.ToLower(env)
}
