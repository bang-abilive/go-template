package pgsql

import (
	"context"
	"fmt"
	"ndinhbang/go-skeleton/pkg/config"
	"time"

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


// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (p *pgsqlConnection) Health(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	if err := p.connection.Ping(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats, err
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// pgx.Conn does not expose pool-level metrics. Report connection-level flags only.
	stats["closed"] = fmt.Sprintf("%t", p.connection.PgConn().IsClosed())
	stats["tx_status"] = string(p.connection.PgConn().TxStatus())

	return stats, nil
}
