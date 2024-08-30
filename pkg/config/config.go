package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ooyeku/grav-lsm/embedded"
)

// Config represents the configuration settings for the application.
// It contains settings for the database, server, and logging.
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

// DatabaseConfig represents the configuration for connecting to a database.
// It contains the driver, host, port, user, password, database name, and SSL mode.
type DatabaseConfig struct {
	Driver        string
	Host          string
	Port          int
	User          string
	Password      string
	Name          string
	SSLMode       string
	ContainerName string
}

// ServerConfig represents the configuration for a server, including the host and port it is running on.
type ServerConfig struct {
	Host string
	Port int
}

// LoggingConfig represents the configuration for logging.
//
// It contains the following fields:
//   - Level: the logging level, which can be "debug", "info", "warn", or "error"
//   - File: the file path where the logs will be written, if specified
type LoggingConfig struct {
	Level string
	File  string
}

// LoadConfig reads the embedded config.json file and parses it into a Config object.
// It returns a pointer to the Config object and an error if any occurs during the process.
// The Config object holds the configuration for the program, including the database, server, and logging configurations.
func LoadConfig() (*Config, error) {
	var cfg Config

	// Try to load from local file first
	localConfig, err := os.ReadFile("config.json")
	if err == nil {
		if err := json.Unmarshal(localConfig, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse local config file: %w", err)
		}
	} else {
		// If local file doesn't exist, load from embedded
		embeddedConfig, err := embedded.EmbeddedFiles.ReadFile("config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded config file: %w", err)
		}
		if err := json.Unmarshal(embeddedConfig, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse embedded config file: %w", err)
		}
	}

	setDefaults(&cfg)
	return &cfg, nil
}

// setDefaults sets default values for the given Config object if any of the fields are empty or zero valued.
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
	if config.Database.ContainerName == "" {
		config.Database.ContainerName = "gravorm-db"
	}
}

// GetConfigPath retrieves the path to the configuration file. It first checks if the
// environment variable "GRAVORM_CONFIG_PATH" is set, and if so, returns its value.
// If the environment variable is not set, the function returns the path "." indicating
// the current directory.
func GetConfigPath() string {
	if configPath := os.Getenv("GRAVORM_CONFIG_PATH"); configPath != "" {
		return configPath
	}
	return "."
}

// SaveConfig saves the given configuration to a file specified by GetConfigPath.
// It creates a new file using os.Create and closes it using defer file.Close().
// It then encodes the config using json.NewEncoder and returns any error encountered.
func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
