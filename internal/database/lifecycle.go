package database

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/ooyeku/grav-orm/pkg/config"
)

// DBLifecycleManager is a struct that manages the lifecycle of the database
type DBLifecycleManager struct {
	// Add fields if necessary
	config *config.Config
}

func NewDBLifecycleManager(config *config.Config) *DBLifecycleManager {
	return &DBLifecycleManager{
		config: config,
	}
}

func (dm *DBLifecycleManager) BuildDatabaseImage() error {
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

func (dm *DBLifecycleManager) StartDatabaseContainer() error {
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

func (dm *DBLifecycleManager) StopDatabaseContainer() error {
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

func (dm *DBLifecycleManager) RemoveDatabaseContainer() error {
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
