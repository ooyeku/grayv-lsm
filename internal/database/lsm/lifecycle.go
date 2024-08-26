package lsm

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/ooyeku/grav-orm/pkg/config"
)

// DBLifecycleManager is responsible for managing the lifecycle of a database.
// It provides methods to build a database image, start a database container, stop a database container,
// and remove a database container. The configuration for the database is provided through the `config`
// field of type `*config.Config`, which contains the necessary parameters such as driver, host, port, username, password, database name, and SSL mode.
type DBLifecycleManager struct {
	// Add fields if necessary
	config *config.Config
}

// NewDBLifecycleManager initializes a new DBLifecycleManager instance with the given configuration.
// It takes a *config.Config as input and returns a pointer to the initialized DBLifecycleManager.
// The config parameter represents the configuration used by the database lifecycle manager.
// Example usage: dbManager := database.NewDBLifecycleManager(cfg)
func NewDBLifecycleManager(config *config.Config) *DBLifecycleManager {
	return &DBLifecycleManager{
		config: config,
	}
}

// BuildDatabaseImage builds the database Docker image by executing a bash script.
// It uses the "/bin/bash" command with the "internal/database/lsm/build.sh" script.
// The standard output and error output of the script are stored in memory buffers for logging purposes.
// If the script execution fails, an error is returned with the detailed command output and error output.
// The command output is also logged regardless of success or failure.
// This method does not take any arguments and returns an error indicating the success or failure of the build process.
func (dm *DBLifecycleManager) BuildDatabaseImage() error {
	cmd := exec.Command("/bin/bash", "internal/database/lsm/build.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command output: %s", out.String())
		log.Printf("Command error output: %s", stderr.String())
		return fmt.Errorf("failed to build database image: %w", err)
	}

	log.Printf("Command output: %s", out.String())
	return nil
}

// StartDatabaseContainer starts the database container by executing the up.sh script.
//
// It returns an error if the command to start the container fails. The error includes
// the output and error output of the command.
func (dm *DBLifecycleManager) StartDatabaseContainer() error {
	cmd := exec.Command("/bin/bash", "internal/database/lsm/up.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command output: %s", out.String())
		log.Printf("Command error output: %s", stderr.String())
		return fmt.Errorf("failed to start database container: %w", err)
	}

	log.Printf("Command output: %s", out.String())
	return nil
}

// StopDatabaseContainer stops the database container by executing the "down.sh" script.
// It returns an error if the execution of the script fails.
func (dm *DBLifecycleManager) StopDatabaseContainer() error {
	cmd := exec.Command("/bin/bash", "internal/database/lsm/down.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command output: %s", out.String())
		log.Printf("Command error output: %s", stderr.String())
		return fmt.Errorf("failed to stop database container: %w", err)
	}

	log.Printf("Command output: %s", out.String())
	return nil
}

// RemoveDatabaseContainer removes the database Docker container by executing the "remove.sh" script.
// It returns an error if the removal process fails. The error message includes the standard output and
// standard error of the command execution.
func (dm *DBLifecycleManager) RemoveDatabaseContainer() error {
	cmd := exec.Command("/bin/bash", "internal/database/lsm/remove.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command output: %s", out.String())
		log.Printf("Command error output: %s", stderr.String())
		return fmt.Errorf("failed to remove database container: %w", err)
	}

	log.Printf("Command output: %s", out.String())
	return nil
}
