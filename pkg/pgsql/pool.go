package pgsql

import (
	"context"
	"fmt"
	"ndinhbang/go-skeleton/pkg/config"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgsqlPool struct {
	pool *pgxpool.Pool
}

func NewPgsqlPool(ctx context.Context, cfg *config.DatabaseConfig) (*pgsqlPool, error) {
	config, err := pgxpool.ParseConfig(cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("[pgsql/pool] failed to parse database config: %w", err)
	}
	config.MinConns = cfg.MinConns               // Default is 0
	config.MaxConns = cfg.MaxConns               // Default is 4
	config.MaxConnLifetime = cfg.MaxConnLifetime // Default is 1 hour
	config.MaxConnIdleTime = cfg.MaxConnIdleTime // Default is 30 minutes
	config.HealthCheckPeriod = 30 * time.Second  // Default is 1 minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		// https://oneuptime.com/blog/post/2026-02-20-go-error-handling-patterns/view
		return nil, fmt.Errorf("[pgsql/pool] failed to create connection pool: %w", err)
	}

	// Verify connection pool is healthy
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("[pgsql/pool] unable to ping database: %w", err)
	}

	return &pgsqlPool{pool: pool}, nil
}

func (p *pgsqlPool) Close(ctx context.Context) error {
	if p.pool == nil {
		return fmt.Errorf("[pgsql/pool] connection pool is already closed or not initialized")
	}
	p.pool.Close()

	return nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (p *pgsqlPool) Health(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	if err := p.pool.Ping(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats, err
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := p.pool.Stat()

	stats["open_connections"] = strconv.Itoa(int(dbStats.TotalConns()))
	stats["in_use"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["idle"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["wait_count"] = strconv.FormatInt(dbStats.EmptyAcquireCount(), 10)
	stats["wait_duration"] = dbStats.EmptyAcquireWaitTime().String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleDestroyCount(), 10)
	stats["max_lifetime_closed"] = strconv.Itoa(int(dbStats.MaxLifetimeDestroyCount()))

	// Evaluate stats to provide a health message
	if dbStats.TotalConns() > p.pool.Config().MaxConns {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.EmptyAcquireCount() > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleDestroyCount() > int64(dbStats.TotalConns())/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeDestroyCount() > int64(dbStats.TotalConns())/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats, nil
}
