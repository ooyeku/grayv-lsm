package orm

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ooyeku/grav-lsm/pkg/config"
)

type Connection struct {
	db *sql.DB
}

func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Connection{db: db}, nil
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func (c *Connection) Ping() error {
	return c.db.Ping()
}

func (c *Connection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c *Connection) GetDB() *sql.DB {
	return c.db
}

func (c *Connection) ListTables() ([]string, error) {
	rows, err := c.db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
