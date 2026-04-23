package pgsql

import (
	"context"
	"fmt"
	"ndinhbang/go-skeleton/pkg/config"

	"github.com/jackc/pgx/v5"
)

type pgsqlConnection struct {
	connection *pgx.Conn
}

func NewPgsqlConnection(ctx context.Context, cfg *config.DatabaseConfig) (*pgsqlConnection, error) {
	conn, err := pgx.Connect(ctx, cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("[pgsql/db] failed to connect to database: %w", err)
	}
	return &pgsqlConnection{connection: conn}, nil
}

func (c *pgsqlConnection) Close(ctx context.Context) error {
	return c.connection.Close(ctx)
}
