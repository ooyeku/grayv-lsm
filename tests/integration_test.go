package tests

// NOTE: To run this test suite, make sure any postgres instances are stopped and removed before running the tests.

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"fmt"
	"github.com/ooyeku/grayv-lsm/cmd"
	"github.com/ooyeku/grayv-lsm/internal/orm"
	"github.com/ooyeku/grayv-lsm/pkg/config"
	log "github.com/ooyeku/grayv-lsm/pkg/logging"
)

func TestMain(m *testing.M) {
	// Setup
	setupTestEnvironment()

	// Run tests
	code := m.Run()

	// Teardown
	teardownTestEnvironment()

	os.Exit(code)
}

func setupTestEnvironment() {
	// Create a temporary directory for test configuration
	tempDir, err := os.MkdirTemp("", "grav-lsm-test")
	if err != nil {
		panic("Failed to create temp directory: " + err.Error())
	}

	// Create a specific file for the configuration
	configFile := filepath.Join(tempDir, "config.json")

	// Set the GRAVORM_CONFIG_PATH environment variable to the specific file
	os.Setenv("GRAVORM_CONFIG_PATH", configFile)

	// Create a test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpassword",
			Name:     "testdb",
			SSLMode:  "disable",
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Logging: config.LoggingConfig{
			Level: "info",
			File:  "test.log",
		},
	}

	// Save the test configuration
	err = config.SaveConfig(cfg)
	if err != nil {
		panic("Failed to save test configuration: " + err.Error())
	}
}

func teardownTestEnvironment() {
	// Get the temporary directory path
	tempDir := os.Getenv("GRAVORM_CONFIG_PATH")

	// Remove the temporary directory and its contents
	err := os.RemoveAll(tempDir)
	if err != nil {
		panic("Failed to remove test configuration directory: " + err.Error())
	}

	// Unset the environment variable
	os.Unsetenv("GRAVORM_CONFIG_PATH")
}

func TestAppLifecycle(t *testing.T) {
	var databaseStarted bool

	// Setup database
	t.Run("DatabaseLifecycle", func(t *testing.T) {
		testDatabaseLifecycle(t)
		databaseStarted = true
	})

	// Run other tests only if database is started
	if databaseStarted {
		t.Run("ModelOperations", testModelOperations)
		t.Run("ORMOperations", testORMOperations)
	}

	// Cleanup after all tests
	defer func() {
		if databaseStarted {
			cmd.RootCmd.SetArgs([]string{"db", "stop"})
			cmd.RootCmd.Execute()
			cmd.RootCmd.SetArgs([]string{"db", "remove"})
			cmd.RootCmd.Execute()
		}
	}()
}

func testDatabaseLifecycle(t *testing.T) {
	// Test database build
	cmd.RootCmd.SetArgs([]string{"db", "build"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to build database: %v", err)
	}

	// Test database start
	cmd.RootCmd.SetArgs([]string{"db", "start"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to start database: %v", err)
	}

	// Wait for the database to be ready
	if err := waitForDatabaseReady(); err != nil {
		t.Fatalf("Database failed to become ready: %v", err)
	}

	// Set up the database schema
	if err := setupDatabaseSchema(); err != nil {
		t.Fatalf("Failed to set up database schema: %v", err)
	}

	// Test database status
	cmd.RootCmd.SetArgs([]string{"db", "status"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to get database status: %v", err)
	}
}

// Rename the duplicate function
func waitForDatabaseReady() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	log.NewColorfulLogger().Info("Waiting for database to become ready...")
	for i := 0; i < 60; i++ { // Increase to 60 attempts (30 seconds)
		conn, err := orm.NewConnection(&cfg.Database)
		if err != nil {
			log.NewColorfulLogger().Infof("Attempt %d: Connection error: %v", i+1, err)
		} else {
			defer conn.Close()
			// Try to create a test table to ensure the database is ready
			_, err := conn.GetDB().Exec("CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY)")
			if err == nil {
				// Drop the test table
				_, _ = conn.GetDB().Exec("DROP TABLE test_table")
				log.NewColorfulLogger().Info("Database is ready")
				return nil
			}
			log.NewColorfulLogger().Infof("Attempt %d: Table creation error: %v", i+1, err)
		}
		log.NewColorfulLogger().Infof("Attempt %d: Database not ready yet, waiting...", i+1)
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("database did not become ready in time")
}

func waitForDatabaseReadyBenchmark() error {
	return waitForDatabaseReady()
}

func setupDatabaseSchema() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	conn, err := orm.NewConnection(&cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Create the models table
	_, err = conn.GetDB().Exec(`
		CREATE TABLE IF NOT EXISTS models (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			fields JSONB NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create models table: %v", err)
	}

	// Create the users table
	_, err = conn.GetDB().Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	return nil
}

func testModelOperations(t *testing.T) {
	// Test model creation
	cmd.RootCmd.SetArgs([]string{"model", "create", "TestModel", "--fields", "name:string,age:int"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Test model list
	cmd.RootCmd.SetArgs([]string{"model", "list"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}

	// Test model update
	cmd.RootCmd.SetArgs([]string{"model", "update", "TestModel", "--add-fields", "email:string"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to update model: %v", err)
	}

	// Test model generation
	cmd.RootCmd.SetArgs([]string{"model", "generate", "TestModel"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to generate model: %v", err)
	}
}

func testORMOperations(t *testing.T) {
	// Test ORM query
	cmd.RootCmd.SetArgs([]string{"orm", "query", "SELECT * FROM users"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute ORM query: %v", err)
	}

}

func setupDatabaseSchemaForBenchmark(b *testing.B) error {
	// Reuse the existing setupDatabaseSchema function
	return setupDatabaseSchema()
}

func BenchmarkDatabaseLifecycle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Test database build
		cmd.RootCmd.SetArgs([]string{"db", "build"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to build database: %v", err)
		}

		// Test database start
		cmd.RootCmd.SetArgs([]string{"db", "start"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to start database: %v", err)
		}

		// Wait for the database to be ready
		if err := waitForDatabaseReadyBenchmark(); err != nil {
			b.Fatalf("Database failed to become ready: %v", err)
		}

		// Set up the database schema
		if err := setupDatabaseSchemaForBenchmark(b); err != nil {
			b.Fatalf("Failed to set up database schema: %v", err)
		}

		// Test database status
		cmd.RootCmd.SetArgs([]string{"db", "status"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to get database status: %v", err)
		}

		// Test database stop
		cmd.RootCmd.SetArgs([]string{"db", "stop"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to stop database: %v", err)
		}

		// Test database remove
		cmd.RootCmd.SetArgs([]string{"db", "remove"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to remove database: %v", err)
		}
	}
}

func BenchmarkDatabaseBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd.RootCmd.SetArgs([]string{"db", "build"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to build database: %v", err)
		}
	}
}

func BenchmarkDatabaseStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd.RootCmd.SetArgs([]string{"db", "start"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to start database: %v", err)
		}
	}
}

func BenchmarkModelCommands(b *testing.B) {
	// setup database
	cmd.RootCmd.SetArgs([]string{"db", "build"})
	if err := cmd.RootCmd.Execute(); err != nil {
		b.Fatalf("Failed to build database: %v", err)
	}

	cmd.RootCmd.SetArgs([]string{"db", "start"})
	if err := cmd.RootCmd.Execute(); err != nil {
		b.Fatalf("Failed to start database: %v", err)
	}
	for i := 0; i < b.N; i++ {
		// Test model creation
		cmd.RootCmd.SetArgs([]string{"model", "create", "TestModel", "--fields", "name:string,age:int"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to create model: %v", err)
		}

		// Test model list
		cmd.RootCmd.SetArgs([]string{"model", "list"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to list models: %v", err)
		}

		// Test model update
		cmd.RootCmd.SetArgs([]string{"model", "update", "TestModel", "--add-fields", "email:string"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to update model: %v", err)
		}

		// Test model generation
		cmd.RootCmd.SetArgs([]string{"model", "generate", "TestModel"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to generate model: %v", err)
		}
	}
}

func BenchmarkDatabaseStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd.RootCmd.SetArgs([]string{"db", "stop"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to stop database: %v", err)
		}
	}
}

func BenchmarkDatabaseRemove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd.RootCmd.SetArgs([]string{"db", "remove"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to remove database: %v", err)
		}
	}
}

func BenchmarkDatabaseStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd.RootCmd.SetArgs([]string{"db", "status"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to get database status: %v", err)
		}
	}
}

func BenchmarkAppCommands(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Test app creation
		cmd.RootCmd.SetArgs([]string{"app", "create", "TestApp"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to create app: %v", err)
		}

		// Test app list
		cmd.RootCmd.SetArgs([]string{"app", "list"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to list apps: %v", err)
		}

		// Test app deletion
		cmd.RootCmd.SetArgs([]string{"app", "delete", "TestApp"})
		if err := cmd.RootCmd.Execute(); err != nil {
			b.Fatalf("Failed to delete app: %v", err)
		}
	}
}
