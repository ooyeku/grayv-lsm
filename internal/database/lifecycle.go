package database

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/ooyeku/grav-orm/pkg/config"
)

// Define the DatabaseManager struct
type DatabaseManager struct {
	// Add fields if necessary
	config *config.Config
}

func NewDatabaseManager(config *config.Config) *DatabaseManager {
	return &DatabaseManager{
		config: config,
	}
}

func (dm *DatabaseManager) BuildDatabaseImage() error {
	cmd := exec.Command("/bin/bash", "internal/database/build.sh")
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

func (dm *DatabaseManager) StartDatabaseContainer() error {
	cmd := exec.Command("/bin/bash", "internal/database/up.sh")
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

func (dm *DatabaseManager) StopDatabaseContainer() error {
	cmd := exec.Command("/bin/bash", "internal/database/down.sh")
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

func (dm *DatabaseManager) RemoveDatabaseContainer() error {
	cmd := exec.Command("/bin/bash", "internal/database/remove.sh")
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
