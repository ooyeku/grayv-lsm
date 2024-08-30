package cmd

import (
	"fmt"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"github.com/ooyeku/grav-lsm/pkg/logging"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var configLogger *logging.ColorfulLogger

func init() {
	configLogger = logging.NewColorfulLogger()
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Grav LSM configuration",
	Long:  `View or edit the Grav LSM configuration settings.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run:   runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run:   runConfigSet,
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	RootCmd.AddCommand(configCmd)
}

func runConfigGet(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		configLogger.Error(fmt.Sprintf("Error loading config: %v", err))
		return
	}

	value := getConfigValue(cfg, args[0])
	if value != "" {
		configLogger.Info(fmt.Sprintf("%s: %s", args[0], value))
	} else {
		configLogger.Warn(fmt.Sprintf("Configuration key '%s' not found", args[0]))
	}
}

func runConfigSet(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		configLogger.Error(fmt.Sprintf("Error loading config: %v", err))
		return
	}

	if setConfigValue(cfg, args[0], args[1]) {
		err = config.SaveConfig(cfg)
		if err != nil {
			configLogger.Error(fmt.Sprintf("Error saving config: %v", err))
			return
		}
		configLogger.Info(fmt.Sprintf("Configuration updated: %s = %s", args[0], args[1]))
	} else {
		configLogger.Warn(fmt.Sprintf("Configuration key '%s' not found", args[0]))
	}
}

func getConfigValue(cfg *config.Config, key string) string {
	switch strings.ToLower(key) {
	case "database.driver":
		return cfg.Database.Driver
	case "database.host":
		return cfg.Database.Host
	case "database.port":
		return fmt.Sprintf("%d", cfg.Database.Port)
	case "database.user":
		return cfg.Database.User
	case "database.password":
		return cfg.Database.Password
	case "database.name":
		return cfg.Database.Name
	case "database.sslmode":
		return cfg.Database.SSLMode
	case "server.host":
		return cfg.Server.Host
	case "server.port":
		return fmt.Sprintf("%d", cfg.Server.Port)
	case "logging.level":
		return cfg.Logging.Level
	case "logging.file":
		return cfg.Logging.File
	case "database.containername":
		return cfg.Database.ContainerName
	default:
		return ""
	}
}

func setConfigValue(cfg *config.Config, key, value string) bool {
	switch strings.ToLower(key) {
	case "database.driver":
		cfg.Database.Driver = value
	case "database.host":
		cfg.Database.Host = value
	case "database.port":
		cfg.Database.Port = parseInt(value)
	case "database.user":
		cfg.Database.User = value
	case "database.password":
		cfg.Database.Password = value
	case "database.name":
		cfg.Database.Name = value
	case "database.sslmode":
		cfg.Database.SSLMode = value
	case "server.host":
		cfg.Server.Host = value
	case "server.port":
		cfg.Server.Port = parseInt(value)
	case "logging.level":
		cfg.Logging.Level = value
	case "logging.file":
		cfg.Logging.File = value
	case "database.containername":
		cfg.Database.ContainerName = value
	default:
		return false
	}
	return true
}

func parseInt(value string) int {
	i, _ := strconv.Atoi(value)
	return i
}
