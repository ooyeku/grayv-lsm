package seed

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// Seed represents a single seed file
type Seed struct {
	Name string
	SQL  string
}

// Seeder manages the database seeding process
type Seeder struct {
	db    *sql.DB
	seeds []*Seed
}

// NewSeeder creates a new Seeder instance
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// LoadSeeds loads all seed files from the specified directory
func (s *Seeder) LoadSeeds(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		logrus.WithError(err).Error("failed to read seeds directory")
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			seed, err := parseSeedFile(filepath.Join(dir, file.Name()))
			if err != nil {
				logrus.WithError(err).Errorf("failed to parse seed file %s", file.Name())
				return err
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

// Seed executes all loaded seed files
func (s *Seeder) Seed() error {
	for _, seed := range s.seeds {
		if err := s.executeSeed(seed); err != nil {
			return err
		}
	}
	return nil
}

// executeSeed runs a single seed file
func (s *Seeder) executeSeed(seed *Seed) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithError(err).Error("error starting transaction")
		return err
	}

	if _, err := tx.Exec(seed.SQL); err != nil {
		tx.Rollback()
		logrus.WithError(err).Errorf("error executing seed %s", seed.Name)
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logrus.WithError(err).Errorf("error committing seed %s", seed.Name)
		return err
	}

	logrus.Infof("Executed seed: %s", seed.Name)
	return nil
}

// parseSeedFile reads and parses a seed file
func parseSeedFile(filename string) (*Seed, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Seed{
		Name: filepath.Base(filename),
		SQL:  strings.TrimSpace(string(content)),
	}, nil
}
