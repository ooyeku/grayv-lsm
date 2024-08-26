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
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
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
		return err
	}

	for _, migration := range m.migrations {
		if err := m.runMigration(migration); err != nil {
			return err
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
		return err
	}

	for i := 0; i < steps && i < len(appliedMigrations); i++ {
		if err := m.rollbackMigration(appliedMigrations[i]); err != nil {
			return err
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
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = $1", migration.Version).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking migration status: %w", err)
	}

	if count > 0 {
		return nil
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	if _, err := tx.Exec(migration.UpSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("error applying migration %s: %w", migration.Name, err)
	}

	if _, err := tx.Exec("INSERT INTO migrations (version, name) VALUES ($1, $2)",
		migration.Version, migration.Name); err != nil {
		tx.Rollback()
		return fmt.Errorf("error recording migration %s: %w", migration.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing migration %s: %w", migration.Name, err)
	}

	fmt.Printf("Applied migration: %s\n", migration.Name)
	return nil
}

func (m *Migrator) rollbackMigration(migration *Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	if _, err := tx.Exec(migration.DownSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("error rolling back migration %s: %w", migration.Name, err)
	}

	if _, err := tx.Exec("DELETE FROM migrations WHERE version = $1", migration.Version); err != nil {
		tx.Rollback()
		return fmt.Errorf("error removing migration record %s: %w", migration.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing rollback of migration %s: %w", migration.Name, err)
	}

	fmt.Printf("Rolled back migration: %s\n", migration.Name)
	return nil
}

func (m *Migrator) getAppliedMigrations() ([]*Migration, error) {
	rows, err := m.db.Query("SELECT version, name FROM migrations ORDER BY version DESC")
	if err != nil {
		return nil, fmt.Errorf("error querying migrations: %w", err)
	}
	defer rows.Close()

	var appliedMigrations []*Migration
	for rows.Next() {
		var version int64
		var name string
		if err := rows.Scan(&version, &name); err != nil {
			return nil, fmt.Errorf("error scanning migration row: %w", err)
		}

		for _, migration := range m.migrations {
			if migration.Version == version {
				appliedMigrations = append(appliedMigrations, migration)
				break
			}
		}
	}

	return appliedMigrations, nil
}

func parseMigrationFile(filename string) (*Migration, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
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
		return nil, err
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
