package db

import (
	"context"
	"database/sql"
	"fmt"
	"ndinhbang/go-template/pkg/config"
)

// Ensure PostgresDatabaseCompat implements the Database interface.
var _ Database = (*PostgresDatabaseCompat)(nil)

type PostgresDatabaseCompat struct {
	db *sql.DB
}

func NewPostgresDatabaseCompat(ctx context.Context, cfg *config.DatabaseConfig) (*PostgresDatabaseCompat, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("[db/compat] open db: %w", err)
	}
	db.SetMaxIdleConns(int(cfg.MaxIdleConns))
	db.SetMaxOpenConns(int(cfg.MaxConns))
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	return &PostgresDatabaseCompat{db: db}, nil
}

func (c *PostgresDatabaseCompat) Close() error {
	return c.db.Close()
}