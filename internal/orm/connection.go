package orm

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ooyeku/grayv-lsm/pkg/config"
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

func (c *Connection) CountTables() (int, error) {
	tables, err := c.ListTables()
	if err != nil {
		return 0, fmt.Errorf("failed to list tables: %w", err)
	}
	return len(tables), nil
}

type DatabaseMetrics struct {
	TableCount        int
	DatabaseSize      string
	ActiveConnections int
	Uptime            string
	Commits           int
	Rollbacks         int
	CacheHitRatio     float64
	SlowQueryCount    int
}

func (c *Connection) GetDatabaseMetrics() (*DatabaseMetrics, error) {
	metrics := &DatabaseMetrics{}

	// Fetch table count
	err := c.db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&metrics.TableCount)
	if err != nil {
		return nil, fmt.Errorf("error counting tables: %w", err)
	}

	// Fetch database size
	err = c.db.QueryRow("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&metrics.DatabaseSize)
	if err != nil {
		return nil, fmt.Errorf("error getting database size: %w", err)
	}

	// Fetch active connections
	err = c.db.QueryRow("SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&metrics.ActiveConnections)
	if err != nil {
		return nil, fmt.Errorf("error counting active connections: %w", err)
	}

	// Fetch uptime
	err = c.db.QueryRow("SELECT now() - pg_postmaster_start_time()").Scan(&metrics.Uptime)
	if err != nil {
		return nil, fmt.Errorf("error getting uptime: %w", err)
	}

	// Fetch transaction statistics
	err = c.db.QueryRow("SELECT xact_commit, xact_rollback FROM pg_stat_database WHERE datname = current_database()").Scan(&metrics.Commits, &metrics.Rollbacks)
	if err != nil {
		return nil, fmt.Errorf("error getting transaction statistics: %w", err)
	}

	// Fetch cache hit ratio
	err = c.db.QueryRow(`
		SELECT 
			CASE 
				WHEN sum(heap_blks_hit) + sum(heap_blks_read) = 0 THEN 0
				ELSE sum(heap_blks_hit) * 100.0 / (sum(heap_blks_hit) + sum(heap_blks_read))
			END
		FROM pg_statio_user_tables
	`).Scan(&metrics.CacheHitRatio)
	if err != nil {
		return nil, fmt.Errorf("error calculating cache hit ratio: %w", err)
	}

	// Fetch slow query count
	err = c.db.QueryRow("SELECT COUNT(*) FROM pg_stat_activity WHERE state = 'active' AND now() - query_start > interval '1 hour'").Scan(&metrics.SlowQueryCount)
	if err != nil {
		return nil, fmt.Errorf("error counting slow queries: %w", err)
	}

	return metrics, nil
}
