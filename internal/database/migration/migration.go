package migration

import (
	"database/sql"
	"fmt"
	"github.com/ooyeku/grayv-lsm/embedded"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// logger is a pointer to a logrus.Logger instance. It is used for logging messages and
// configuring log settings.
var logger *logrus.Logger

// init initializes the logger variable by assigning a new instance of logrus.Logger to it.
func init() {
	logger = logrus.New()
}

// Migration represents a database migration.
//
// A Migration struct contains the following fields:
//   - Version: int64 - the version number of the migration
//   - Name: string - the name of the migration
//   - UpSQL: string - the SQL code to apply the migration
//   - DownSQL: string - the SQL code to rollback the migration
//   - Timestamp: time.Time - the timestamp when the migration was created
type Migration struct {
	Version   int64
	Name      string
	UpSQL     string
	DownSQL   string
	Timestamp time.Time
}

// Migrator represents a database migrator that can apply and rollback migrations.
// It keeps track of applied migrations and provides methods for running and undoing
// migrations.
//
// Fields:
// - db: The *sql.DB instance representing the database connection.
// - migrations: A slice of *Migration instances representing the available migrations.
// - logger: The *logrus.Logger instance used for logging migration events.
//
// Usage:
// - To create a new Migrator instance, use the NewMigrator function.
// - To load migrations from the embedded files, use the LoadMigrations method.
// - To apply all available migrations that haven't been applied yet, use the Migrate method.
// - To rollback a specific number of applied migrations, use the Rollback method.
//
// Example usage:
//
//	db, _ := sql.Open("postgres", "postgres://user:pass@localhost/db")
//	logger := logrus.New()
//	migrator := NewMigrator(db, logger)
//	migrator.LoadMigrations()
//	err := migrator.Migrate()
//	err = migrator.Rollback(1)
type Migrator struct {
	db         *sql.DB
	migrations []*Migration
	logger     *logrus.Logger
}

// NewMigrator creates a new instance of Migrator.
// It accepts a *sql.DB database connection and a *logrus.Logger logger.
// Returns a pointer to Migrator struct.
// Example usage:
//
//	migrator := migration.NewMigrator(conn.GetDB(), log)
func NewMigrator(db *sql.DB, logger *logrus.Logger) *Migrator {
	return &Migrator{db: db, logger: logger}
}

// LoadMigrations reads and loads the embedded migration files from the "migrations" directory.
// It reads the files with the ".sql" extension,
// parses each migration file,
// sorts the migrations based on their version,
// and appends them to the Migrator's migrations slice.
// Returns an error if there is any issue reading, parsing, or sorting the migrations.
func (m *Migrator) LoadMigrations() error {
	entries, err := embedded.EmbeddedFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations directory: %w", err)
	}

	var loadErrors []error
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sql" {
			migrationContent, err := embedded.EmbeddedFiles.ReadFile(filepath.Join("migrations", entry.Name()))
			if err != nil {
				loadErrors = append(loadErrors, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err))
				continue
			}
			migration, err := parseMigrationContent(entry.Name(), string(migrationContent))
			if err != nil {
				loadErrors = append(loadErrors, fmt.Errorf("failed to parse migration file %s: %w", entry.Name(), err))
				continue
			}
			m.migrations = append(m.migrations, migration)
		}
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	if len(loadErrors) > 0 {
		return fmt.Errorf("errors occurred while loading migrations: %v", loadErrors)
	}

	return nil
}

// parseMigrationContent parses the content of a migration file and returns a *Migration object
// containing the parsed information. The function splits the content into two parts, using "-- Down"
// as the delimiter. If the content does not have exactly two parts, it returns an error. It then trims
// the whitespace from both parts and assigns them to the UpSQL and DownSQL fields of the *Migration object.
// It also calls parseVersionFromFilename to parse the version from the given filename. If there is an error
// parsing the version, it returns an error. Finally, it initializes a new *Migration object with the parsed
// information, including the version, filename, timestamp (set to the current time), and returns it along
// with nil error.
func parseMigrationContent(filename, content string) (*Migration, error) {
	parts := strings.Split(content, "-- Down")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid migration file format")
	}

	upSQL := strings.TrimSpace(parts[0])
	downSQL := strings.TrimSpace(parts[1])

	version, err := parseVersionFromFilename(filename)
	if err != nil {
		return nil, fmt.Errorf("error parsing version from filename: %w", err)
	}

	return &Migration{
		Version:   version,
		Name:      filename,
		UpSQL:     upSQL,
		DownSQL:   downSQL,
		Timestamp: time.Now(),
	}, nil
}

// Migrate applies pending migrations to the database.
// It creates the migrations table if it does not exist.
// It retrieves the list of applied migrations from the database.
// For each migration that has not been applied, it runs the migration.
// Returns an error if any step fails.
func (m *Migrator) Migrate() error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range m.migrations {
		if !contains(appliedMigrations, migration.Version) {
			if err := m.runMigration(migration); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
			}
		}
	}

	return nil
}

// Rollback rolls back a specified number of migrations by executing their corresponding down SQL statements.
// It retrieves the list of applied migrations, finds the migration to be rolled back,
// and then executes the rollback process by running the migration's down SQL statement.
// The steps parameter determines the number of migrations to roll back.
// If steps is less than or equal to 0, the function returns immediately without performing any rollback operations.
// If there are fewer applied migrations than the specified steps, it only rolls back the available migrations.
// The function returns an error if it encounters any issues during the rollback process.
func (m *Migrator) Rollback(steps int) error {
	if steps <= 0 {
		return nil
	}

	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for i := 0; i < steps && i < len(appliedMigrations); i++ {
		migration := m.findMigration(appliedMigrations[i])
		if migration == nil {
			return fmt.Errorf("migration with version %d not found", appliedMigrations[i])
		}
		if err := m.rollbackMigration(migration); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Name, err)
		}
	}

	return nil
}

const migrationsTableName = "migrations"

// createMigrationsTable creates a table called "migrations" in the database if it does not exist already.
// The table has three columns: "version" of type BIGINT and primary key, "name" of type TEXT and not null,
// and "applied_at" of type TIMESTAMP WITH TIME ZONE with a default value of the current timestamp.
// This method returns an error if there was a problem executing the SQL statement to create the table.
func (m *Migrator) createMigrationsTable() error {
	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            version BIGINT PRIMARY KEY,
            name TEXT NOT NULL,
            applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `, migrationsTableName)
	_, err := m.db.Exec(query)
	return err
}

// runMigration applies a migration to the database using a transaction.
// It executes the UpSQL statement of the migration and inserts a record
// of the migration into the migrations table.
// If an error occurs at any step, the transaction is rolled back.
//
// Parameters:
// - migration: The migration to be applied.
//
// Returns:
// - error: An error if any occurred during the migration process.
func (m *Migrator) runMigration(migration *Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("error applying migration: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO migrations (version, name) VALUES ($1, $2)",
		migration.Version, migration.Name); err != nil {
		return fmt.Errorf("error recording migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing migration: %w", err)
	}

	m.logger.Infof("Applied migration: %s", migration.Name)
	return nil
}

// rollbackMigration rolls back a migration by executing the DownSQL statement and removing the migration record from the database.
// It starts a transaction, rolls it back in case of an error, and commits the rollback if successful.
// It logs the name of the rolled-back migration.
// It returns an error if any operation fails.
func (m *Migrator) rollbackMigration(migration *Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.DownSQL); err != nil {
		return fmt.Errorf("error rolling back migration: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM migrations WHERE version = $1", migration.Version); err != nil {
		return fmt.Errorf("error removing migration record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing rollback: %w", err)
	}

	m.logger.Infof("Rolled back migration: %s", migration.Name)
	return nil
}

// getAppliedMigrations queries the migrations table in the database and retrieves
// the versions of the applied migrations, ordered in descending order. It returns
// a slice of int64 representing the versions and an error if there was any issue
// querying the database.
func (m *Migrator) getAppliedMigrations() ([]int64, error) {
	rows, err := m.db.Query("SELECT version FROM migrations ORDER BY version DESC")
	if err != nil {
		return nil, fmt.Errorf("error querying migrations: %w", err)
	}
	defer rows.Close()

	var appliedMigrations []int64
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("error scanning migration row: %w", err)
		}
		appliedMigrations = append(appliedMigrations, version)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over migration rows: %w", err)
	}

	return appliedMigrations, nil
}

// findMigration searches for a migration with the specified version in the list of migrations.
// It returns a pointer to the found migration, or nil if no migration with that version was found.
func (m *Migrator) findMigration(version int64) *Migration {
	for _, migration := range m.migrations {
		if migration.Version == version {
			return migration
		}
	}
	return nil
}

// parseVersionFromFilename extracts the version number from a migration filename.
// It splits the filename by '_' and checks if there are at least two parts.
// If the version part cannot be converted to an int64, it returns an error.
// Returns the parsed version number as an int64 and nil or an error if the format is invalid.
func parseVersionFromFilename(filename string) (int64, error) {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid migration filename format")
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid migration version: %w", err)
	}

	return version, nil
}

// contains checks if an item is present in a slice of int64 values.
// It iterates through the slice and returns true if the item is found,
// otherwise it returns false.
// The function takes two parameters:
// - slice: the slice of int64 values to be searched
// - item: the item to be checked if it is present in the slice
// It returns a boolean value indicating whether the item is present in the slice or not.
func contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
