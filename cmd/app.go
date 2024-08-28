package cmd

import (
	"github.com/ooyeku/grav-lsm/internal/app"
	"github.com/spf13/cobra"
)

// appCreator is a pointer variable of type *app.AppCreator. It represents an application creator that can be used
// to create, list, and delete Grav apps. The appCreator variable should be initialized with an instance of *app.AppCreator
// before it can be used. Once initialized, it can be used to call methods like CreateApp, ListApps, and DeleteApp to perform
// the corresponding operations. The appCreator variable can also be used as a dependency for other parts of the application
// that may need to create or manipulate Grav apps.
var appCreator *app.AppCreator

// appCmd is a variable of type *cobra.Command that represents a command for managing Grav apps.
// It has the following properties:
//   - Use:   "app" (the string used to invoke the command)
//   - Short: "Manage Grav apps" (a short description of the command)
//   - Long:  "Create, list, and delete Grav apps." (a long description of the command)
//
// This variable is used in the init() function to register subcommands and add the command to the RootCmd.
// Example usage:
//
//	appCreator = app.NewAppCreator()
//	appCmd.AddCommand(createAppCmd)
//	appCmd.AddCommand(listAppsCmd)
//	appCmd.AddCommand(deleteAppCmd)
//	RootCmd.AddCommand(appCmd)
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage Grav apps",
	Long:  `Create, list, and delete Grav apps.`,
}

// createAppCmd represents the Cobra command for creating a new Grav app. It takes one argument, the name of the app,
// and uses the AppCreator to create the app. It logs successful creation and any errors encountered during the process.
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

// listAppsCmd is a variable of type *cobra.Command that represents the "list" command.
// It is used to list all Grav apps. The command defines a Run function that calls the ListApps method
// of the appCreator instance to get a list of Grav apps. It then logs the apps or an appropriate message.
// If any error occurs during the command execution, an error message is logged.
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

// deleteAppCmd is a variable of type *cobra.Command that represents the command for deleting a Grav app.
// It is used as a subcommand of the appCmd command. When executed, it calls the DeleteApp method on appCreator
// to delete the app with the specified name. If the deletion is successful, the app name is logged using the
// logger's Info method. If the deletion fails, an error is logged using the logger's Errorf method.
// The app name is passed as an argument to the command.
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

// init is a function that initializes the application creator and adds
// commands to the root command. It adds the create, list, and delete app
// commands to the appCmd. Then, it adds the appCmd to the rootCmd.
func init() {
	appCreator = app.NewAppCreator()

	appCmd.AddCommand(createAppCmd)
	appCmd.AddCommand(listAppsCmd)
	appCmd.AddCommand(deleteAppCmd)
	RootCmd.AddCommand(appCmd)
}
