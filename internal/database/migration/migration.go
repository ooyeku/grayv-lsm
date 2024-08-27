package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Migration struct {
	Version   int64
	Name      string
	UpSQL     string
	DownSQL   string
	Timestamp time.Time
}

type Migrator struct {
	db         *sql.DB
	migrations []*Migration
	logger     *logrus.Logger
}

func NewMigrator(db *sql.DB, logger *logrus.Logger) *Migrator {
	return &Migrator{db: db, logger: logger}
}

func (m *Migrator) LoadMigrations(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migration, err := parseMigrationFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to parse migration file %s: %w", file.Name(), err)
			}
			m.migrations = append(m.migrations, migration)
		}
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	return nil
}

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

func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

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

	return appliedMigrations, nil
}

func (m *Migrator) findMigration(version int64) *Migration {
	for _, migration := range m.migrations {
		if migration.Version == version {
			return migration
		}
	}
	return nil
}

func parseMigrationFile(filename string) (*Migration, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading migration file: %w", err)
	}

	parts := strings.Split(string(content), "-- Down")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid migration file format")
	}

	upSQL := strings.TrimSpace(parts[0])
	downSQL := strings.TrimSpace(parts[1])

	name := filepath.Base(filename)
	version, err := parseVersionFromFilename(name)
	if err != nil {
		return nil, fmt.Errorf("error parsing version from filename: %w", err)
	}

	return &Migration{
		Version:   version,
		Name:      name,
		UpSQL:     upSQL,
		DownSQL:   downSQL,
		Timestamp: time.Now(),
	}, nil
}

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

func contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
