package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"ndinhbang/go-skeleton/pkg/config"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // register database/sql driver name "pgx"
)

type pgsqlCompat struct {
	db *sql.DB
}

func NewPgsqlCompat(ctx context.Context, cfg *config.DatabaseConfig) (*pgsqlCompat, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("[pgsql/compat] failed to open database: %w", err)
	}
	db.SetMaxIdleConns(int(cfg.MaxIdleConns))
	db.SetMaxOpenConns(int(cfg.MaxConns))
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	return &pgsqlCompat{db: db}, nil
}

func (c *pgsqlCompat) Close(ctx context.Context) error {
	return c.db.Close()
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (c *pgsqlCompat) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := c.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := c.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}
