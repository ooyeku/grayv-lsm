package cmd

import (
	"fmt"
	"github.com/ooyeku/grav-orm/internal/database"
	"github.com/ooyeku/grav-orm/pkg/config"
	"github.com/spf13/cobra"
	"os/exec"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the database lifecycle",
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the database Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(config.GetConfigPath())
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		dbManager := database.NewDBLifecycleManager(cfg)
		if err := dbManager.BuildDatabaseImage(); err != nil {
			fmt.Println("Error building database image:", err)
		} else {
			fmt.Println("Database image built successfully")
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(config.GetConfigPath())
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		dbManager := database.NewDBLifecycleManager(cfg)
		if err := dbManager.StartDatabaseContainer(); err != nil {
			fmt.Println("Error starting database container:", err)
		} else {
			fmt.Println("Database container started successfully")
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(config.GetConfigPath())
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		dbManager := database.NewDBLifecycleManager(cfg)
		if err := dbManager.StopDatabaseContainer(); err != nil {
			fmt.Println("Error stopping database container:", err)
		} else {
			fmt.Println("Database container stopped successfully")
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(config.GetConfigPath())
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		dbManager := database.NewDBLifecycleManager(cfg)
		if err := dbManager.RemoveDatabaseContainer(); err != nil {
			fmt.Println("Error removing database container:", err)
		} else {
			fmt.Println("Database container removed successfully")
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if the database Docker container is running",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("docker", "ps", "-f", "name=gravorm-db", "--format", "{{.Names}}").Output()
		if err != nil {
			fmt.Println("Error checking database container status:", err)
			return
		}
		if string(out) == "gravorm-db\n" {
			fmt.Println("Database container is running")
		} else {
			fmt.Println("Database container is not running")
		}
	},
}

func init() {
	dbCmd.AddCommand(buildCmd)
	dbCmd.AddCommand(startCmd)
	dbCmd.AddCommand(stopCmd)
	dbCmd.AddCommand(removeCmd)
	dbCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(dbCmd)
}
