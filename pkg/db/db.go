package db

import (
	"context"
	"fmt"
	"ndinhbang/go-template/pkg/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Close() error
}

// Ensure PostgresDatabase implements the Database interface.
var _ Database = (*PostgresDatabase)(nil)

type PostgresDatabase struct {
	pool *pgxpool.Pool
}

func NewPostgresDatabase(ctx context.Context, cfg *config.DatabaseConfig) (*PostgresDatabase, error) {
	dbcfg, err := pgxpool.ParseConfig(cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("[db] parse db config: %w", err)
	}
	dbcfg.MinConns = cfg.MinConns               // Default is 0
	dbcfg.MaxConns = cfg.MaxConns               // Default is 4
	dbcfg.MaxConnLifetime = cfg.MaxConnLifetime // Default is 1 hour
	dbcfg.MaxConnIdleTime = cfg.MaxConnIdleTime // Default is 30 minutes
	dbcfg.HealthCheckPeriod = 30 * time.Second  // Default is 1 minute

	pool, err := pgxpool.NewWithConfig(ctx, dbcfg)
	if err != nil {
		// https://oneuptime.com/blog/post/2026-02-20-go-error-handling-patterns/view
		return nil, fmt.Errorf("[db] create pool : %w", err)
	}
	// fail fast if the database is not reachable or credentials are wrong
	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	// Verify connection pool is healthy
	if err := pool.Ping(ctxPing); err != nil {
		return nil, fmt.Errorf("[db] ping: %w", err)
	}

	return &PostgresDatabase{pool: pool}, nil
}

func (db *PostgresDatabase) Close() error {
	db.pool.Close()
	return nil
}

func (db PostgresDatabase) GetPool() *pgxpool.Pool {
	return db.pool
}
