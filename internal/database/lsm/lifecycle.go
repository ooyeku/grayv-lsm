package lsm

import (
	"fmt"
	"github.com/ooyeku/grav-lsm/pkg/config"
	"os"
	"os/exec"
	"strings"
)

type DBLifecycleManager struct {
	config *config.Config
}

func NewDBLifecycleManager(config *config.Config) *DBLifecycleManager {
	return &DBLifecycleManager{
		config: config,
	}
}

func (dm *DBLifecycleManager) setEnvVars() {
	os.Setenv("DB_USER", dm.config.Database.User)
	os.Setenv("DB_PASSWORD", dm.config.Database.Password)
	os.Setenv("DB_NAME", dm.config.Database.Name)
}

func (dm *DBLifecycleManager) fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (dm *DBLifecycleManager) runCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (dm *DBLifecycleManager) BuildImage() error {
	dm.setEnvVars()
	fmt.Println("Starting the build process for the database Docker image...")

	if !dm.fileExists("./internal/database/lsm/Dockerfile") {
		return fmt.Errorf("Dockerfile not found!")
	}

	buildCommand := "docker build -f ./internal/database/lsm/Dockerfile -t gravorm-db --build-arg DB_USER=$DB_USER --build-arg DB_PASSWORD=$DB_PASSWORD --build-arg DB_NAME=$DB_NAME ."
	output, err := dm.runCommand(buildCommand)
	if err != nil {
		return fmt.Errorf("Failed to build the database Docker image: %v\nOutput: %s", err, output)
	}

	fmt.Println("Database Docker image built successfully.")
	return nil
}

func (dm *DBLifecycleManager) StartContainer() error {
	dm.setEnvVars()
	fmt.Println("Starting the database Docker container...")

	// Check if the container already exists
	output, _ := dm.runCommand("docker ps -aq -f name=gravorm-db")
	if output != "" {
		fmt.Println("Container gravorm-db already exists. Removing it...")
		_, err := dm.runCommand("docker rm -f gravorm-db")
		if err != nil {
			return fmt.Errorf("Failed to remove existing container: %v", err)
		}
	}

	// Check if the image exists locally
	output, _ = dm.runCommand("docker images -q gravorm-db")
	if output == "" {
		return fmt.Errorf("Docker image gravorm-db not found. Please build the image first.")
	}

	// Start the Docker container
	startCommand := fmt.Sprintf("docker run -d --name gravorm-db -e POSTGRES_USER=%s -e POSTGRES_PASSWORD=%s -e POSTGRES_DB=%s -p 5432:5432 gravorm-db",
		dm.config.Database.User, dm.config.Database.Password, dm.config.Database.Name)
	output, err := dm.runCommand(startCommand)
	if err != nil {
		return fmt.Errorf("Failed to start the database Docker container: %v\nOutput: %s", err, output)
	}

	fmt.Println("Database Docker container started successfully.")

	// Verify the container is running
	output, err = dm.runCommand("docker ps -q -f name=gravorm-db")
	if err != nil || output == "" {
		return fmt.Errorf("Database Docker container is not running.")
	}

	// Verify environment variables inside the container
	output, err = dm.runCommand("docker exec gravorm-db env | grep POSTGRES")
	if err != nil {
		return fmt.Errorf("Failed to verify environment variables in the container: %v\nOutput: %s", err, output)
	}

	fmt.Println("Environment variables are set correctly in the container.")
	return nil
}

func (dm *DBLifecycleManager) StopContainer() error {
	fmt.Println("Stopping the database Docker container...")
	output, err := dm.runCommand("docker stop gravorm-db")
	if err != nil {
		return fmt.Errorf("Failed to stop the database Docker container: %v\nOutput: %s", err, output)
	}
	fmt.Println("Database Docker container stopped successfully.")
	return nil
}

func (dm *DBLifecycleManager) RemoveContainer() error {
	fmt.Println("Removing the database Docker container...")
	output, err := dm.runCommand("docker rm gravorm-db")
	if err != nil {
		return fmt.Errorf("Failed to remove the database Docker container: %v\nOutput: %s", err, output)
	}
	fmt.Println("Database Docker container removed successfully.")
	return nil
}

func (dm *DBLifecycleManager) GetStatus() (string, error) {
	// Check if the container exists
	output, err := dm.runCommand("docker ps -a --filter name=gravorm-db --format '{{.Status}}'")
	if err != nil {
		return "", fmt.Errorf("Failed to get the status of the database Docker container: %v", err)
	}

	output = strings.TrimSpace(output)
	if output == "" {
		return "Container does not exist", nil
	}

	// Check if the container is running
	isRunning, err := dm.runCommand("docker inspect -f '{{.State.Running}}' gravorm-db")
	if err != nil {
		return "", fmt.Errorf("Failed to inspect the database Docker container: %v", err)
	}

	isRunning = strings.TrimSpace(isRunning)
	if isRunning == "true" {
		return fmt.Sprintf("Container is running. Status: %s", output), nil
	} else {
		return fmt.Sprintf("Container is not running. Status: %s", output), nil
	}
}
