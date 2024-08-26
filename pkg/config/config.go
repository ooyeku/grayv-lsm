package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
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
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	v.SetEnvPrefix("GRAVORM")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct: %w", err)
	}

	// Set defaults if not provided
	setDefaults(&config)

	return &config, nil
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
