package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("GRAVORM_CONFIG_PATH", "test_config.json")
	config := &Config{
		Database: DatabaseConfig{
			Driver:   "test",
			Host:     "test",
			Port:     1000,
			User:     "test",
			Password: "test",
			Name:     "test",
			SSLMode:  "test",
		},
		Server: ServerConfig{
			Host: "test",
			Port: 1000,
		},
		Logging: LoggingConfig{
			Level: "test",
			File:  "test",
		},
	}
	SaveConfig(config)

	_, err := LoadConfig()
	if err != nil {
		t.Fatalf("wanted nil but got %v", err)
	}

	os.Remove("test_config.json")
}

func TestSaveConfig(t *testing.T) {
	os.Setenv("GRAVORM_CONFIG_PATH", "test_config.json")
	config := &Config{
		Database: DatabaseConfig{
			Driver:   "test",
			Host:     "test",
			Port:     1000,
			User:     "test",
			Password: "test",
			Name:     "test",
			SSLMode:  "test",
		},
		Server: ServerConfig{
			Host: "test",
			Port: 1000,
		},
		Logging: LoggingConfig{
			Level: "test",
			File:  "test",
		},
	}
	err := SaveConfig(config)
	if err != nil {
		t.Fatalf("wanted nil but got %v", err)
	}

	_, err = LoadConfig()
	if err != nil {
		t.Fatalf("wanted nil but got %v", err)
	}

	var resultingConfig Config
	file, _ := os.Open("test_config.json")
	json.NewDecoder(file).Decode(&resultingConfig)

	if config.Database.Driver != resultingConfig.Database.Driver ||
		config.Server.Host != resultingConfig.Server.Host ||
		config.Logging.Level != resultingConfig.Logging.Level {
		t.Fatalf("Incorrect config file created")
	}

	file.Close()
	os.Remove("test_config.json")
}

func TestSetDefaults(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{},
		Server:   ServerConfig{},
		Logging:  LoggingConfig{},
	}

	setDefaults(config)

	if config.Database.Driver != "postgres" ||
		config.Database.Host != "localhost" ||
		config.Database.Port != 5432 ||
		config.Database.SSLMode != "disable" ||
		config.Server.Host != "0.0.0.0" ||
		config.Server.Port != 8080 ||
		config.Logging.Level != "info" {
		t.Fatalf("Default config not set correctly")
	}
}
