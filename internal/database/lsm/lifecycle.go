package lsm

import (
	"context"
	"fmt"
	"github.com/ooyeku/grayv-lsm/embedded"
	"github.com/ooyeku/grayv-lsm/pkg/config"
	"github.com/ooyeku/grayv-lsm/pkg/logging"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// log is a variable of type logrus.Logger. It is used for logging messages and errors throughout the program.
var log *logging.ColorfulLogger

// init initializes the logging configuration for the application.
//
// It sets the formatter to use a text formatter with full timestamps,
// sets the output to standard output, and sets the log level to Info.
//
// This function is called automatically when the package is initialized.
func init() {
	log = logging.NewColorfulLogger()
}

// DBLifecycleManager represents a type that manages the lifecycle of a database. It contains a config.Config object that
// holds the configuration for the program. The DBLifecycleManager is responsible for setting environment variables,
// checking file existence, running commands, building and starting a Docker container, stopping and removing the container, and
// getting the status of the container.
type DBLifecycleManager struct {
	config        *config.Config
	logger        *logging.ColorfulLogger
	containerName string
}

// NewDBLifecycleManager creates a new instance of the DBLifecycleManager struct.
// It takes a pointer to a config.Config object as a parameter and returns a pointer to the newly created DBLifecycleManager object.
func NewDBLifecycleManager(cfg *config.Config) *DBLifecycleManager {
	return &DBLifecycleManager{
		config:        cfg,
		logger:        logging.NewColorfulLogger(),
		containerName: cfg.Database.ContainerName,
	}
}

// Acknowledge that setEnvVars is intentionally unused
var _ = (*DBLifecycleManager).setEnvVars

func (dm *DBLifecycleManager) setEnvVars() error {
	vars := map[string]string{
		"DB_USER":           dm.config.Database.User,
		"DB_PASSWORD":       dm.config.Database.Password,
		"DB_NAME":           dm.config.Database.Name,
		"DB_CONTAINER_NAME": dm.config.Database.ContainerName,
		"DB_IMAGE":          dm.config.Database.Image,
	}

	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}

	return nil
}

// Acknowledge that fileExists is intentionally unused
var _ = (*DBLifecycleManager).fileExists

func (dm *DBLifecycleManager) fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Update the runCommand method signature
func (dm *DBLifecycleManager) runCommand(command string, args ...interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(command, args...))
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command timed out")
	}
	return string(output), err
}

// BuildImage builds the Docker image for the database using the specified Dockerfile.
// It sets the necessary environment variables, checks if the Dockerfile exists,
// and runs the build command. If the build process fails, it logs the error and returns it.
// Otherwise, it logs the successful build and returns nil.
func (dm *DBLifecycleManager) BuildImage() error {
	dockerfileContent, err := embedded.EmbeddedFiles.ReadFile("Dockerfile")
	if err != nil {
		return fmt.Errorf("failed to read embedded Dockerfile: %w", err)
	}

	// Remove the COPY instruction
	dockerfileLines := strings.Split(string(dockerfileContent), "\n")
	var newDockerfileContent strings.Builder
	for _, line := range dockerfileLines {
		if !strings.HasPrefix(line, "COPY ./internal/database/init.sql") {
			newDockerfileContent.WriteString(line + "\n")
		}
	}

	tempDir, err := os.MkdirTemp("", "grayv-db-build")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.WithError(err).Error("failed to remove temp directory")
		}
	}()

	if err := os.WriteFile(filepath.Join(tempDir, "Dockerfile"), []byte(newDockerfileContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile to temp directory: %w", err)
	}

	buildCommand := fmt.Sprintf("docker build -t %s %s", dm.config.Database.Image, tempDir)
	output, err := dm.runCommand(buildCommand)
	if err != nil {
		return fmt.Errorf("failed to build the database docker image: %v\nOutput: %s", err, output)
	}

	log.Infof("Database Docker image %s built successfully.", dm.config.Database.Image)
	return nil
}

// StartContainer starts the database Docker container.
// It checks if the container already exists and removes it if it does.
// It checks if the image exists locally and returns an error if it does not.
// It starts the Docker container by running a command with necessary environment variables.
// It verifies that the container is running and that the environment variables are set correctly inside the container.
// Returns an error if any step fails.
func (dm *DBLifecycleManager) StartContainer() error {
	log.Infof("Starting the database Docker container %s...", dm.config.Database.ContainerName)

	// Check if the container already exists
	output, _ := dm.runCommand(fmt.Sprintf("docker ps -aq -f name=%s", dm.config.Database.ContainerName))
	if output != "" {
		log.Infof("Container %s already exists. Removing it...", dm.config.Database.ContainerName)
		_, err := dm.runCommand(fmt.Sprintf("docker rm -f %s", dm.config.Database.ContainerName))
		if err != nil {
			return fmt.Errorf("failed to remove existing container: %v", err)
		}
	}

	// Check if the image exists locally
	output, _ = dm.runCommand(fmt.Sprintf("docker images -q %s", dm.config.Database.Image))
	if output == "" {
		return fmt.Errorf("docker image %s not found. Please build the image first", dm.config.Database.Image)
	}

	// Start the Docker container
	startCommand := fmt.Sprintf("docker run -d --name %s -e POSTGRES_USER=%s -e POSTGRES_PASSWORD=%s -e POSTGRES_DB=%s -p 5432:5432 %s",
		dm.config.Database.ContainerName, dm.config.Database.User, dm.config.Database.Password, dm.config.Database.Name, dm.config.Database.Image)
	output, err := dm.runCommand(startCommand)
	if err != nil {
		return fmt.Errorf("failed to start the database docker container: %v\nOutput: %s", err, output)
	}

	log.Infof("Database Docker container %s started successfully.", dm.config.Database.ContainerName)

	// Verify the container is running
	output, err = dm.runCommand(fmt.Sprintf("docker ps -q -f name=%s", dm.config.Database.ContainerName))
	if err != nil || output == "" {
		return fmt.Errorf("database Docker container is not running")
	}

	// Verify environment variables inside the container
	output, err = dm.runCommand(fmt.Sprintf("docker exec %s env | grep POSTGRES", dm.config.Database.ContainerName))
	if err != nil {
		return fmt.Errorf("failed to verify environment variables in the container: %v\nOutput: %s", err, output)
	}

	log.Infof("Environment variables are set correctly in the container %s.", dm.config.Database.ContainerName)
	return nil
}

// StopContainer stops the database Docker container by running the command "docker stop gravorm-db".
// It returns an error if it fails to stop the container, along with the output of the command.
// If the container is stopped successfully, it logs a success message and returns nil.
func (dm *DBLifecycleManager) StopContainer() error {
	log.Infof("Stopping the database Docker container %s...", dm.containerName)
	output, err := dm.runCommand(fmt.Sprintf("docker stop %s", dm.containerName))
	if err != nil {
		return fmt.Errorf("failed to stop the database Docker container: %v\nOutput: %s", err, output)
	}
	log.Infof("Database Docker container %s stopped successfully.", dm.containerName)
	return nil
}

// RemoveContainer removes the database Docker container. It runs the "docker rm gravorm-db" command
// to remove the container. If the command fails, it returns an error with the failure message.
// Otherwise, it logs a success message and returns nil.
func (dm *DBLifecycleManager) RemoveContainer() error {
	log.Infof("Removing the database Docker container %s...", dm.config.Database.ContainerName)
	output, err := dm.runCommand(fmt.Sprintf("docker rm %s", dm.config.Database.ContainerName))
	if err != nil {
		return fmt.Errorf("failed to remove the database Docker container: %v\nOutput: %s", err, output)
	}
	log.Infof("Database Docker container %s removed successfully.", dm.config.Database.ContainerName)
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
	output, err := dm.runCommand(fmt.Sprintf("docker ps -a --filter name=%s --format '{{.Status}}'", dm.config.Database.ContainerName))
	if err != nil {
		log.WithError(err).Error("failed to get the status of the database Docker container")
		return "", fmt.Errorf("failed to get the status of the database Docker container: %v", err)
	}

	output = strings.TrimSpace(output)
	if output == "" {
		log.Infof("Container %s does not exist", dm.config.Database.ContainerName)
		return fmt.Sprintf("container %s does not exist", dm.config.Database.ContainerName), nil
	}

	// Check if the container is running
	isRunning, err := dm.runCommand(fmt.Sprintf("docker inspect -f '{{.State.Running}}' %s", dm.config.Database.ContainerName))
	if err != nil {
		log.WithError(err).Error("failed to inspect the database Docker container")
		return "", fmt.Errorf("failed to inspect the database Docker container: %v", err)
	}

	isRunning = strings.TrimSpace(isRunning)
	if isRunning == "true" {
		status := fmt.Sprintf("Container %s is running. Status: %s", dm.config.Database.ContainerName, output)
		log.Info(status)
		return status, nil
	} else {
		status := fmt.Sprintf("Container %s is not running. Status: %s", dm.config.Database.ContainerName, output)
		log.Info(status)
		return status, nil
	}
}
