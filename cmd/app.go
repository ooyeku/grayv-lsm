/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ooyeku/grav-lsm/internal/app"
	"github.com/spf13/cobra"
)

var appCreator *app.AppCreator

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage Grav apps",
	Long:  `Create, list, and delete Grav apps.`,
}

var createAppCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new Grav app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		if err := appCreator.CreateApp(appName); err != nil {
			log.WithError(err).Errorf("Failed to create Grav app '%s'", appName)
		} else {
			log.Infof("Grav app '%s' created successfully", appName)
		}
	},
}

var listAppsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Grav apps",
	Run: func(cmd *cobra.Command, args []string) {
		apps, err := appCreator.ListApps()
		if err != nil {
			log.WithError(err).Error("Failed to list Grav apps")
			return
		}
		if len(apps) == 0 {
			log.Info("No Grav apps found")
		} else {
			log.Info("Grav apps:")
			for _, app := range apps {
				log.Infof("- %s", app)
			}
		}
	},
}

var deleteAppCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a Grav app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		if err := appCreator.DeleteApp(appName); err != nil {
			log.WithError(err).Errorf("Failed to delete Grav app '%s'", appName)
		} else {
			log.Infof("Grav app '%s' deleted successfully", appName)
		}
	},
}

func init() {
	appCreator = app.NewAppCreator()

	appCmd.AddCommand(createAppCmd)
	appCmd.AddCommand(listAppsCmd)
	appCmd.AddCommand(deleteAppCmd)
	RootCmd.AddCommand(appCmd)
}
