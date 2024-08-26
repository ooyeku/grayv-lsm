package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuin/gopher-lua"
	"os/exec"
)

func runLuaFunction(funcName string) error {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile("internal/database/lsm/db_lsm.lua"); err != nil {
		return err
	}
	fn := L.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		return fmt.Errorf("'%s' is not a function", funcName)
	}
	if err := L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}); err != nil {
		return err
	}
	return nil
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the database lifecycle",
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the database Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLuaFunction("Build_image"); err != nil {
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
		if err := runLuaFunction("Start_container"); err != nil {
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
		if err := runLuaFunction("Stop_container"); err != nil {
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
		if err := runLuaFunction("Remove_container"); err != nil {
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
