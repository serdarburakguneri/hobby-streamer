package config

import (
	"os"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	config := manager.GetConfig()
	if config.Service == "" {
		t.Error("Service name should not be empty")
	}

	if config.Environment == "" {
		t.Error("Environment should not be empty")
	}

	if config.Server.Port == "" {
		t.Error("Server port should not be empty")
	}
}

func TestEnvironmentVariableOverride(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("LOG_LEVEL")
	}()

	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	config := manager.GetConfig()
	if config.Server.Port != "9090" {
		t.Errorf("Expected port '9090', got '%s'", config.Server.Port)
	}

	if config.Log.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.Log.Level)
	}
}

func TestSecretsManager(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "test-key")
	os.Setenv("NEO4J_PASSWORD", "test-password")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("NEO4J_PASSWORD")
	}()

	secretsManager := NewSecretsManager()
	secretsManager.LoadFromEnvironment()

	if secretsManager.Get("aws_access_key_id") != "test-key" {
		t.Error("Failed to get AWS access key from secrets")
	}

	if secretsManager.Get("neo4j_password") != "test-password" {
		t.Error("Failed to get Neo4j password from secrets")
	}

	if secretsManager.GetOrDefault("nonexistent", "default") != "default" {
		t.Error("GetOrDefault should return default value for nonexistent key")
	}
}

func TestFeatureFlagsManager(t *testing.T) {
	featureManager := NewFeatureFlagsManager()

	featureManager.Register("test_feature", "Test feature description", false)
	if featureManager.IsEnabled("test_feature") {
		t.Error("Feature should be disabled by default")
	}

	featureManager.Enable("test_feature")
	if !featureManager.IsEnabled("test_feature") {
		t.Error("Feature should be enabled after calling Enable")
	}

	featureManager.Disable("test_feature")
	if featureManager.IsEnabled("test_feature") {
		t.Error("Feature should be disabled after calling Disable")
	}

	featureManager.SetPercentage("test_feature", 50)
	flag := featureManager.Get("test_feature")
	if flag.Percentage != 50 {
		t.Errorf("Expected percentage 50, got %d", flag.Percentage)
	}
}

func TestConfigurationValidation(t *testing.T) {
	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	config := manager.GetConfig()

	if config.Service == "" {
		t.Error("Service name should be validated and not empty")
	}

	if config.Environment == "" {
		t.Error("Environment should be validated and not empty")
	}

	if config.Server.Port == "" {
		t.Error("Server port should be validated and not empty")
	}
}

func TestDefaultValues(t *testing.T) {
	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	config := manager.GetConfig()

	if config.Server.ReadTimeout != 15*time.Second {
		t.Errorf("Expected read timeout 15s, got %v", config.Server.ReadTimeout)
	}

	if config.Server.WriteTimeout != 15*time.Second {
		t.Errorf("Expected write timeout 15s, got %v", config.Server.WriteTimeout)
	}

	if config.Server.IdleTimeout != 60*time.Second {
		t.Errorf("Expected idle timeout 60s, got %v", config.Server.IdleTimeout)
	}
}

func TestDynamicConfig(t *testing.T) {
	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	dynamicCfg := manager.GetDynamicConfig()

	if dynamicCfg == nil {
		t.Error("Dynamic config should not be nil")
	}

	if dynamicCfg.GetBaseConfig() == nil {
		t.Error("Base config should not be nil")
	}

	if dynamicCfg.HasComponent("nonexistent") {
		t.Error("Non-existent component should return false")
	}

	if dynamicCfg.GetComponent("nonexistent") != nil {
		t.Error("Non-existent component should return nil")
	}
}

func TestDynamicConfigAccessors(t *testing.T) {
	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	dynamicCfg := manager.GetDynamicConfig()

	if dynamicCfg.GetStringFromComponent("nonexistent", "key") != "" {
		t.Error("GetStringFromComponent should return empty string for non-existent component")
	}

	if dynamicCfg.GetIntFromComponent("nonexistent", "key") != 0 {
		t.Error("GetIntFromComponent should return 0 for non-existent component")
	}

	if dynamicCfg.GetBoolFromComponent("nonexistent", "key") {
		t.Error("GetBoolFromComponent should return false for non-existent component")
	}

	if dynamicCfg.GetFloatFromComponent("nonexistent", "key") != 0.0 {
		t.Error("GetFloatFromComponent should return 0.0 for non-existent component")
	}

	defaultDuration := 30 * time.Minute
	if dynamicCfg.GetDurationFromComponent("nonexistent", "key", defaultDuration) != defaultDuration {
		t.Error("GetDurationFromComponent should return default value for non-existent component")
	}

	if dynamicCfg.GetComponentAsMap("nonexistent") != nil {
		t.Error("GetComponentAsMap should return nil for non-existent component")
	}

	if dynamicCfg.GetComponentAsStringMap("nonexistent") != nil {
		t.Error("GetComponentAsStringMap should return nil for non-existent component")
	}
}

func TestDynamicConfigWithComponents(t *testing.T) {
	manager, err := NewManager("test-service")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	dynamicCfg := manager.GetDynamicConfig()

	if dynamicCfg.GetComponent("aws") != nil {
		t.Error("AWS component should be nil when not configured")
	}

	if dynamicCfg.GetComponent("neo4j") != nil {
		t.Error("Neo4j component should be nil when not configured")
	}

	if dynamicCfg.GetComponent("redis") != nil {
		t.Error("Redis component should be nil when not configured")
	}

	if dynamicCfg.GetComponent("keycloak") != nil {
		t.Error("Keycloak component should be nil when not configured")
	}
}
