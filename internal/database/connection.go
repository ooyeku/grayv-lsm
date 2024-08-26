package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ooyeku/grav-orm/pkg/config"
)

type Connection struct {
	db *sql.DB
}

// NewConnection creates a new database connection using the provided configuration
func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Connection{db: db}, nil
}
