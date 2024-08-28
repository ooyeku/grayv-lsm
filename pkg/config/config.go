package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ooyeku/grav-lsm/embedded"
)

// Config holds all configuration for our program
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

// DatabaseConfig holds all database-related configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig holds all server-related configuration
type ServerConfig struct {
	Host string
	Port int
}

// LoggingConfig holds all logging-related configuration
type LoggingConfig struct {
	Level string
	File  string
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	configData, err := embedded.EmbeddedFiles.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(config *Config) {
	if config.Database.Driver == "" {
		config.Database.Driver = "postgres"
	}
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	if configPath := os.Getenv("GRAVORM_CONFIG_PATH"); configPath != "" {
		return configPath
	}
	return "."
}

func SaveConfig(cfg *Config) error {
	file, err := os.Create(GetConfigPath())
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(cfg)
}
