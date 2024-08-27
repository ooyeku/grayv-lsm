package lsm

import (
	"fmt"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// log is a variable of type logrus.Logger. It is used for logging messages and errors throughout the program.
var log = logrus.New()

// init initializes the logging configuration for the application.
//
// It sets the formatter to use a text formatter with full timestamps,
// sets the output to standard output, and sets the log level to Info.
//
// This function is called automatically when the package is initialized.
func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

// DBLifecycleManager represents a type that manages the lifecycle of a database. It contains a config.Config object that
// holds the configuration for the program. The DBLifecycleManager is responsible for setting environment variables,
// checking file existence, running commands, building and starting a Docker container, stopping and removing the container, and
// getting the status of the container.
type DBLifecycleManager struct {
	config *config.Config
}

// NewDBLifecycleManager creates a new instance of the DBLifecycleManager struct.
// It takes a pointer to a config.Config object as a parameter and returns a pointer to the newly created DBLifecycleManager object.
func NewDBLifecycleManager(config *config.Config) *DBLifecycleManager {
	return &DBLifecycleManager{
		config: config,
	}
}

// setEnvVars sets the environment variables for the database connection. It uses the values from the `config.Database`
// field of the `DBLifecycleManager` instance to set the `DB_USER`, `DB_PASSWORD`, and `DB_NAME` environment variables.
// If setting any of these variables fails, an error is logged and the method returns without further action.
func (dm *DBLifecycleManager) setEnvVars() {
	err := os.Setenv("DB_USER", dm.config.Database.User)
	if err != nil {
		log.WithError(err).Error("failed to set environment variable DB_USER")
		return
	}
	err = os.Setenv("DB_PASSWORD", dm.config.Database.Password)
	if err != nil {
		log.WithError(err).Error("failed to set environment variable DB_PASSWORD")
		return
	}
	err = os.Setenv("DB_NAME", dm.config.Database.Name)
	if err != nil {
		log.WithError(err).Error("failed to set environment variable DB_NAME")
		return
	}
}

// fileExists checks if a file exists in the filesystem.
// It takes a name string parameter representing the file path or name.
// It returns a bool value indicating whether the file exists or not.
// The function uses the os.Stat function to check the file's information,
// and the os.IsNotExist function to determine if the file doesn't exist.
// This function is used internally by other methods in the DBLifecycleManager struct.
// Example usage:
//
//	dm.fileExists("./internal/database/lsm/Dockerfile")
func (dm *DBLifecycleManager) fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// runCommand runs a shell command and returns the combined output and any error that occurred.
func (dm *DBLifecycleManager) runCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// BuildImage builds the Docker image for the database using the specified Dockerfile.
// It sets the necessary environment variables, checks if the Dockerfile exists,
// and runs the build command. If the build process fails, it logs the error and returns it.
// Otherwise, it logs the successful build and returns nil.
func (dm *DBLifecycleManager) BuildImage() error {
	dm.setEnvVars()
	log.Info("Starting the build process for the database Docker image...")

	if !dm.fileExists("./internal/database/lsm/Dockerfile") {
		log.Error("Dockerfile not found!")
		return fmt.Errorf("dockerfile not found")
	}

	buildCommand := "docker build -f ./internal/database/lsm/Dockerfile -t gravorm-db --build-arg DB_USER=$DB_USER --build-arg DB_PASSWORD=$DB_PASSWORD --build-arg DB_NAME=$DB_NAME ."
	output, err := dm.runCommand(buildCommand)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":  err,
			"output": output,
		}).Error("failed to build the database Docker image")
		return fmt.Errorf("failed to build the database Docker image: %v\nOutput: %s", err, output)
	}

	log.Info("Database Docker image built successfully.")
	return nil
}

// StartContainer starts the database Docker container.
// It checks if the container already exists and removes it if it does.
// It checks if the image exists locally and returns an error if it does not.
// It starts the Docker container by running a command with necessary environment variables.
// It verifies that the container is running and that the environment variables are set correctly inside the container.
// Returns an error if any step fails.
func (dm *DBLifecycleManager) StartContainer() error {
	dm.setEnvVars()
	log.Info("Starting the database Docker container...")

	// Check if the container already exists
	output, _ := dm.runCommand("docker ps -aq -f name=gravorm-db")
	if output != "" {
		log.Info("Container gravorm-db already exists. Removing it...")
		_, err := dm.runCommand("docker rm -f gravorm-db")
		if err != nil {
			return fmt.Errorf("failed to remove existing container: %v", err)
		}
	}

	// Check if the image exists locally
	output, _ = dm.runCommand("docker images -q gravorm-db")
	if output == "" {
		return fmt.Errorf("docker image gravorm-db not found. Please build the image first")
	}

	// Start the Docker container
	startCommand := fmt.Sprintf("docker run -d --name gravorm-db -e POSTGRES_USER=%s -e POSTGRES_PASSWORD=%s -e POSTGRES_DB=%s -p 5432:5432 gravorm-db",
		dm.config.Database.User, dm.config.Database.Password, dm.config.Database.Name)
	output, err := dm.runCommand(startCommand)
	if err != nil {
		return fmt.Errorf("failed to start the database docker container: %v\nOutput: %s", err, output)
	}

	log.Info("Database Docker container started successfully.")

	// Verify the container is running
	output, err = dm.runCommand("docker ps -q -f name=gravorm-db")
	if err != nil || output == "" {
		return fmt.Errorf("database Docker container is not running.")
	}

	// Verify environment variables inside the container
	output, err = dm.runCommand("docker exec gravorm-db env | grep POSTGRES")
	if err != nil {
		return fmt.Errorf("failed to verify environment variables in the container: %v\nOutput: %s", err, output)
	}

	log.Info("Environment variables are set correctly in the container.")
	return nil
}

// StopContainer stops the database Docker container by running the command "docker stop gravorm-db".
// It returns an error if it fails to stop the container, along with the output of the command.
// If the container is stopped successfully, it logs a success message and returns nil.
func (dm *DBLifecycleManager) StopContainer() error {
	log.Info("Stopping the database Docker container...")
	output, err := dm.runCommand("docker stop gravorm-db")
	if err != nil {
		return fmt.Errorf("failed to stop the database Docker container: %v\nOutput: %s", err, output)
	}
	log.Info("Database Docker container stopped successfully.")
	return nil
}

// RemoveContainer removes the database Docker container. It runs the "docker rm gravorm-db" command
// to remove the container. If the command fails, it returns an error with the failure message.
// Otherwise, it logs a success message and returns nil.
func (dm *DBLifecycleManager) RemoveContainer() error {
	log.Info("Removing the database Docker container...")
	output, err := dm.runCommand("docker rm gravorm-db")
	if err != nil {
		return fmt.Errorf("failed to remove the database Docker container: %v\nOutput: %s", err, output)
	}
	log.Info("Database Docker container removed successfully.")
	return nil
}

// GetStatus returns the status of the database Docker container.
// It checks if the container exists and if it is running.
// If the container does not exist, it returns "container does not exist".
// If the container is running, it returns "Container is running. Status: <status>".
// If the container is not running, it returns "Container is not running. Status: <status>".
// It returns an error if there is any failure in getting the status of the container.
// The function uses Docker CLI commands to check the status.
func (dm *DBLifecycleManager) GetStatus() (string, error) {
	// Check if the container exists
	output, err := dm.runCommand("docker ps -a --filter name=gravorm-db --format '{{.Status}}'")
	if err != nil {
		log.WithError(err).Error("failed to get the status of the database Docker container")
		return "", fmt.Errorf("failed to get the status of the database Docker container: %v", err)
	}

	output = strings.TrimSpace(output)
	if output == "" {
		log.Info("container does not exist")
		return "container does not exist", nil
	}

	// Check if the container is running
	isRunning, err := dm.runCommand("docker inspect -f '{{.State.Running}}' gravorm-db")
	if err != nil {
		log.WithError(err).Error("failed to inspect the database Docker container")
		return "", fmt.Errorf("failed to inspect the database Docker container: %v", err)
	}

	isRunning = strings.TrimSpace(isRunning)
	if isRunning == "true" {
		status := fmt.Sprintf("Container is running. Status: %s", output)
		log.Info(status)
		return status, nil
	} else {
		status := fmt.Sprintf("Container is not running. Status: %s", output)
		log.Info(status)
		return status, nil
	}
}
