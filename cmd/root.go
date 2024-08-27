package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "grav-lsm",
	Short: "Grav LSM (Lifecycle Management)",
	Long: `
	grav-lsm is a CLI tool for managing the lifecycle of Grav App.  Grav apps are
	lightweight backend components consising of a containerized database, a model/schema generator, and
	an orm system.`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
