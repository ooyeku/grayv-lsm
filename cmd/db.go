package cmd

import (
	"fmt"
	"github.com/ooyeku/grav-lsm/internal/database/lsm"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"github.com/spf13/cobra"
)

var dbManager *lsm.DBLifecycleManager

func init() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	dbManager = lsm.NewDBLifecycleManager(cfg)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the database lifecycle",
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the database Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dbManager.BuildImage(); err != nil {
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
		if err := dbManager.StartContainer(); err != nil {
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
		if err := dbManager.StopContainer(); err != nil {
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
		if err := dbManager.RemoveContainer(); err != nil {
			fmt.Println("Error removing database container:", err)
		} else {
			fmt.Println("Database container removed successfully")
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the health of the database Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := dbManager.GetStatus()
		if err != nil {
			fmt.Println("Error checking database status:", err)
			return
		}
		fmt.Println(status)
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
