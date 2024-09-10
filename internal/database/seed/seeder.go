package seed

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ooyeku/grayv-lsm/embedded"
	"github.com/sirupsen/logrus"
)

// Seed represents a database seed, which encapsulates the name and the SQL statements
// to be executed.
type Seed struct {
	Name string
	SQL  string
}

// Seeder represents a struct for managing database seeding operations.
//
// It contains a database connection (db) and a set of seed objects (seeds).
type Seeder struct {
	db    *sql.DB
	seeds []*Seed
}

// NewSeeder creates a new instance of the Seeder struct which is used to seed the database with initial data.
// It takes a pointer to a sql.DB object as a parameter and returns a pointer to the Seeder struct.
// The sql.DB object is used to execute the SQL queries to seed the database.
// Example usage: seeder := seed.NewSeeder(conn.GetDB())
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// LoadSeeds loads the seed files from the embedded "seeds" directory and populates the Seeder's seeds slice.
// Seed files must have a .sql extension. The seeds are sorted in alphabetical order by filename.
// Returns an error if the embedded seeds directory cannot be read or if any seed file fails to be read.
// This method is part of the Seeder type.
func (s *Seeder) LoadSeeds() error {
	entries, err := embedded.EmbeddedFiles.ReadDir("seeds")
	if err != nil {
		return fmt.Errorf("failed to read embedded seeds directory: %w", err)
	}

	var loadErrors []error
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sql" {
			seedContent, err := embedded.EmbeddedFiles.ReadFile(filepath.Join("seeds", entry.Name()))
			if err != nil {
				loadErrors = append(loadErrors, fmt.Errorf("failed to read seed file %s: %w", entry.Name(), err))
				continue
			}
			seed := &Seed{
				Name: entry.Name(),
				SQL:  string(seedContent),
			}
			s.seeds = append(s.seeds, seed)
		}
	}

	sort.Slice(s.seeds, func(i, j int) bool {
		return s.seeds[i].Name < s.seeds[j].Name
	})

	if len(loadErrors) > 0 {
		return fmt.Errorf("errors occurred while loading seeds: %v", loadErrors)
	}

	return nil
}

// Seed executes all the loaded seeds in the Seeder. Returns an error if any seed fails to execute.
func (s *Seeder) Seed() error {
	for _, seed := range s.seeds {
		if err := s.executeSeed(seed); err != nil {
			return err
		}
	}
	return nil
}

// executeSeed executes the given seed by starting a transaction, executing the SQL statements,
// and committing the transaction. If any error occurs during the process, the transaction
// will be rolled back and the error will be returned. Otherwise, a log message will be printed
// indicating the successful execution of the seed.
//
// Parameters:
// - seed: The seed to be executed.
//
// Returns:
// - An error if any error occurs during the execution of the seed, otherwise nil.
func (s *Seeder) executeSeed(seed *Seed) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithError(err).Error("error starting transaction")
		return err
	}
	defer tx.Rollback()

	// Split the SQL into individual statements
	statements := strings.Split(seed.SQL, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := tx.Exec(stmt); err != nil {
			logrus.WithError(err).Errorf("error executing seed %s", seed.Name)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Errorf("error committing seed %s", seed.Name)
		return err
	}

	logrus.Infof("Executed seed: %s", seed.Name)
	return nil
}
