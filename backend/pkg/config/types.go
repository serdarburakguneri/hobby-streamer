package config

import (
	"time"
)

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
	Test        Environment = "test"
)

type LogConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=text json"`
	Async  struct {
		Enabled    bool `mapstructure:"enabled"`
		BufferSize int  `mapstructure:"buffer_size" validate:"min=1"`
	} `mapstructure:"async"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port" validate:"required"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type SecurityConfig struct {
	RateLimit struct {
		Requests int           `mapstructure:"requests" validate:"min=1"`
		Window   time.Duration `mapstructure:"window" validate:"min=1s"`
	} `mapstructure:"rate_limit"`
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowed_origins"`
		AllowedMethods []string `mapstructure:"allowed_methods"`
		AllowedHeaders []string `mapstructure:"allowed_headers"`
	} `mapstructure:"cors"`
	MaxRequestSize int64 `mapstructure:"max_request_size" validate:"min=1"`
}

type FeatureFlags struct {
	EnableCircuitBreaker bool `mapstructure:"enable_circuit_breaker"`
	EnableRetry          bool `mapstructure:"enable_retry"`
	EnableCaching        bool `mapstructure:"enable_caching"`
	EnableMetrics        bool `mapstructure:"enable_metrics"`
	EnableTracing        bool `mapstructure:"enable_tracing"`
}

type CircuitBreakerConfig struct {
	Threshold int           `mapstructure:"threshold" validate:"min=1"`
	Timeout   time.Duration `mapstructure:"timeout" validate:"min=1s"`
}

type RetryConfig struct {
	MaxAttempts int           `mapstructure:"max_attempts" validate:"min=1"`
	BaseDelay   time.Duration `mapstructure:"base_delay" validate:"min=1ms"`
	MaxDelay    time.Duration `mapstructure:"max_delay" validate:"min=1ms"`
}

type CacheConfig struct {
	TTL time.Duration `mapstructure:"ttl" validate:"min=1s"`
}

type BaseConfig struct {
	Environment    Environment            `mapstructure:"environment" validate:"required,oneof=development staging production test"`
	Service        string                 `mapstructure:"service" validate:"required"`
	Log            LogConfig              `mapstructure:"log"`
	Server         ServerConfig           `mapstructure:"server"`
	Security       SecurityConfig         `mapstructure:"security"`
	Features       FeatureFlags           `mapstructure:"features"`
	CircuitBreaker CircuitBreakerConfig   `mapstructure:"circuit_breaker"`
	Retry          RetryConfig            `mapstructure:"retry"`
	Cache          CacheConfig            `mapstructure:"cache"`
	Components     map[string]interface{} `mapstructure:"components"`
}

type ServiceConfig interface {
	GetBaseConfig() *BaseConfig
	GetComponent(name string) interface{}
	HasComponent(name string) bool
}

type DynamicConfig struct {
	base *BaseConfig
}

func NewDynamicConfig(base *BaseConfig) *DynamicConfig {
	return &DynamicConfig{base: base}
}

func (dc *DynamicConfig) GetBaseConfig() *BaseConfig {
	return dc.base
}

func (dc *DynamicConfig) GetComponent(name string) interface{} {
	if component, exists := dc.base.Components[name]; exists {
		return component
	}
	return nil
}

func (dc *DynamicConfig) HasComponent(name string) bool {
	_, exists := dc.base.Components[name]
	return exists
}

func (dc *DynamicConfig) GetComponentAsMap(name string) map[string]interface{} {
	if component := dc.GetComponent(name); component != nil {
		if componentMap, ok := component.(map[string]interface{}); ok {
			return componentMap
		}
	}
	return nil
}

func (dc *DynamicConfig) GetComponentAsStringMap(name string) map[string]string {
	if componentMap := dc.GetComponentAsMap(name); componentMap != nil {
		result := make(map[string]string)
		for key, value := range componentMap {
			if strValue, ok := value.(string); ok {
				result[key] = strValue
			}
		}
		return result
	}
	return nil
}

func (dc *DynamicConfig) GetStringFromComponent(componentName, key string) string {
	if componentMap := dc.GetComponentAsMap(componentName); componentMap != nil {
		if value, exists := componentMap[key]; exists {
			if strValue, ok := value.(string); ok {
				return strValue
			}
		}
	}
	return ""
}

func (dc *DynamicConfig) GetIntFromComponent(componentName, key string) int {
	if componentMap := dc.GetComponentAsMap(componentName); componentMap != nil {
		if value, exists := componentMap[key]; exists {
			switch v := value.(type) {
			case int:
				return v
			case float64:
				return int(v)
			}
		}
	}
	return 0
}

func (dc *DynamicConfig) GetBoolFromComponent(componentName, key string) bool {
	if componentMap := dc.GetComponentAsMap(componentName); componentMap != nil {
		if value, exists := componentMap[key]; exists {
			if boolValue, ok := value.(bool); ok {
				return boolValue
			}
		}
	}
	return false
}

func (dc *DynamicConfig) GetFloatFromComponent(componentName, key string) float64 {
	if componentMap := dc.GetComponentAsMap(componentName); componentMap != nil {
		if value, exists := componentMap[key]; exists {
			if floatValue, ok := value.(float64); ok {
				return floatValue
			}
		}
	}
	return 0.0
}

func (dc *DynamicConfig) GetDurationFromComponent(componentName, key string, defaultValue time.Duration) time.Duration {
	if componentMap := dc.GetComponentAsMap(componentName); componentMap != nil {
		if value, exists := componentMap[key]; exists {
			if strValue, ok := value.(string); ok {
				if duration, err := time.ParseDuration(strValue); err == nil {
					return duration
				}
			}
		}
	}
	return defaultValue
}

// ---------- Typed component helpers ----------
// KeycloakComponent wraps the generic component map for type-safe access.
// Expected keys: url, realm, client_id, client_secret

type KeycloakComponent struct{ data map[string]interface{} }

func (kc KeycloakComponent) URL() string          { return str(kc.data["url"]) }
func (kc KeycloakComponent) Realm() string        { return str(kc.data["realm"]) }
func (kc KeycloakComponent) ClientID() string     { return str(kc.data["client_id"]) }
func (kc KeycloakComponent) ClientSecret() string { return str(kc.data["client_secret"]) }

// AssetManagerComponent – currently only url field

type AssetManagerComponent struct{ data map[string]interface{} }

func (am AssetManagerComponent) URL() string { return str(am.data["url"]) }

// KafkaComponent example – bootstrap_servers, max_message_bytes

type KafkaComponent struct{ data map[string]interface{} }

func (k KafkaComponent) BootstrapServers() string { return str(k.data["bootstrap_servers"]) }
func (k KafkaComponent) MaxMessageBytes() int     { return intNum(k.data["max_message_bytes"]) }

// Helpers to safely cast
func str(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
func intNum(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	default:
		return 0
	}
}

// Exposed helpers on ComponentConfig
func (dc *DynamicConfig) Keycloak() KeycloakComponent {
	return KeycloakComponent{data: dc.GetComponentAsMap("keycloak")}
}
func (dc *DynamicConfig) AssetManager() AssetManagerComponent {
	return AssetManagerComponent{data: dc.GetComponentAsMap("asset_manager")}
}
func (dc *DynamicConfig) Kafka() KafkaComponent {
	return KafkaComponent{data: dc.GetComponentAsMap("kafka")}
}
