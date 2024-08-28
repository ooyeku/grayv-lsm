package seed

import (
	"database/sql"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ooyeku/grav-lsm/embedded"
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
		logrus.WithError(err).Error("failed to read embedded seeds directory")
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sql" {
			seedContent, err := embedded.EmbeddedFiles.ReadFile(filepath.Join("seeds", entry.Name()))
			if err != nil {
				logrus.WithError(err).Errorf("failed to read seed file %s", entry.Name())
				return err
			}
			seed := &Seed{
				Name: entry.Name(),
				SQL:  string(seedContent),
			}
			s.seeds = append(s.seeds, seed)
		}
	}

	// Sort seeds by filename to ensure consistent order
	sort.Slice(s.seeds, func(i, j int) bool {
		return s.seeds[i].Name < s.seeds[j].Name
	})

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

	if _, err := tx.Exec(seed.SQL); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		logrus.WithError(err).Errorf("error executing seed %s", seed.Name)
		return err
	}

	if err := tx.Commit(); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		logrus.WithError(err).Errorf("error committing seed %s", seed.Name)
		return err
	}

	logrus.Infof("Executed seed: %s", seed.Name)
	return nil
}

// parseSeedFile reads the contents of a seed file specified by the filename parameter,
// and returns a Seed object containing the base filename as the Name property, and the trimmed
// contents of the file as the SQL property. If an error occurs during file reading, the function
// returns nil and the error.
func parseSeedFile(filename string) (*Seed, error) {
	content, err := embedded.EmbeddedFiles.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Seed{
		Name: filepath.Base(filename),
		SQL:  strings.TrimSpace(string(content)),
	}, nil
}
