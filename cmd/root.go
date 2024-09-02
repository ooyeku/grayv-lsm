package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "grayv-lsm",
	Short: "Grayv LSM (Lifecycle Management)",
	Long:  ` grayv-lsm is a CLI tool for managing the lifecycle of Grayv App.  Grayv apps are lightweight backend components consising of a containerized database, a model/schema generator, and an orm system.`,
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
