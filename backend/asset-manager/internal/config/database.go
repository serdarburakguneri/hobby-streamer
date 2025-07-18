package config

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type DatabaseConfig struct {
	Driver neo4j.Driver
}

func NewDatabaseConfig(configManager *config.Manager, secretsManager *config.SecretsManager, log *logger.Logger) (*DatabaseConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	uri := dynamicCfg.GetStringFromComponent("neo4j", "uri")
	username := dynamicCfg.GetStringFromComponent("neo4j", "username")
	password := secretsManager.Get("neo4j_password")

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.MaxConnectionPoolSize = 50
		config.ConnectionAcquisitionTimeout = 30 * time.Second
		config.MaxConnectionLifetime = 1 * time.Hour
		config.ConnectionLivenessCheckTimeout = 30 * time.Second
		config.SocketConnectTimeout = 10 * time.Second
		config.SocketKeepalive = true
	})
	if err != nil {
		log.WithError(err).Error("Failed to create Neo4j driver")
		return nil, err
	}

	if err := driver.VerifyConnectivity(); err != nil {
		log.WithError(err).Error("Failed to connect to Neo4j")
		return nil, err
	}

	log.Info("Neo4j connection established", "uri", uri)

	return &DatabaseConfig{
		Driver: driver,
	}, nil
}

func (dc *DatabaseConfig) Close() {
	if dc.Driver != nil {
		if err := dc.Driver.Close(); err != nil {
			logger.Get().WithError(err).Error("Failed to close Neo4j driver")
		}
	}
}
