package cmd

import (
	"github.com/ooyeku/grayv-lsm/internal/app"
	"github.com/spf13/cobra"
)

var appCreator *app.AppCreator

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage Grayv apps",
	Long:  `Create, list, and delete Grayv apps.`,
}

var createAppCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new Grayv app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		if err := appCreator.CreateApp(appName); err != nil {
			log.WithError(err).Errorf("Failed to create Grayv app '%s'", appName)
		} else {
			log.Infof("Grayv app '%s' created successfully", appName)
		}
	},
}

// listAppsCmd is a variable of type *cobra.Command that represents the "list" command.
// It is used to list all Grav apps. The command defines a Run function that calls the ListApps method
// of the appCreator instance to get a list of Grav apps. It then logs the apps or an appropriate message.
// If any error occurs during the command execution, an error message is logged.
var listAppsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Grayv apps",
	Run: func(cmd *cobra.Command, args []string) {
		apps, err := appCreator.ListApps()
		if err != nil {
			log.WithError(err).Error("Failed to list Grayv apps")
			return
		}
		if len(apps) == 0 {
			log.Info("No Grayv apps found")
		} else {
			log.Info("Grayv apps:")
			for _, app := range apps {
				log.Infof("- %s", app)
			}
		}
	},
}

var deleteAppCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a Grayv app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		if err := appCreator.DeleteApp(appName); err != nil {
			log.WithError(err).Errorf("Failed to delete Grayv app '%s'", appName)
		} else {
			log.Infof("Grayv app '%s' deleted successfully", appName)
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
