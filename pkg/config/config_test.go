package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("GRAVORM_CONFIG_PATH", "config.json")
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

	os.Remove("config.json")
}

func TestSaveConfig(t *testing.T) {
	os.Setenv("GRAVORM_CONFIG_PATH", "config.json")
	config := &Config{
		Database: DatabaseConfig{
			Driver:        "test",
			Host:          "test",
			Port:          1000,
			User:          "test",
			Password:      "test",
			Name:          "test",
			SSLMode:       "test",
			ContainerName: "test-container",
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

	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("wanted nil but got %v", err)
	}

	if !reflect.DeepEqual(config, loadedConfig) {
		t.Errorf("Saved and loaded configs do not match")
		t.Errorf("Original: %+v", config)
		t.Errorf("Loaded: %+v", loadedConfig)
	}

	os.Remove("config.json")
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
