package pgsql

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// *pgsqlCompat must implement [database.Database] (see ndinhbang/go-template/internal/database).
var _ database.Database = (*pgsqlNative)(nil)

// assertUnreachableDBError accepts multiple failure modes when nothing is listening
// on the target (e.g. connection refused, i/o timeout, or context deadline from a short Ping timeout).
func assertUnreachableDBError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	if errors.Is(err, context.DeadlineExceeded) {
		return
	}
	if errors.Is(err, context.Canceled) {
		return
	}
	if opErr, ok := errors.AsType[*net.OpError](err); ok {
		if opErr.Timeout() || opErr.Op == "dial" || opErr.Op == "read" {
			return
		}
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return
	}
	low := strings.ToLower(err.Error())
	for _, sub := range []string{
		"refused",
		"connection reset",
		"reset by peer",
		"dial",
		"connect",
		"i/o timeout",
		"timeout",
		"deadline exceeded",
		"context deadline",
		"no connection",
		"broken pipe",
	} {
		if strings.Contains(low, sub) {
			return
		}
	}
	t.Fatalf("unexpected error (want dial/refused/timeout/deadline): %v", err)
}

func testDatabaseConfig() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		Name:            "testdb",
		Schema:          "public",
		Host:            "127.0.0.1",
		Port:            1,
		User:            "u",
		Password:        "p",
		SSLMode:         "disable",
		MaxConns:        4,
		MinConns:        0,
		MaxConnIdleTime: 10 * time.Minute,
		MaxConnLifetime: time.Hour,
	}
}

func TestNewPgsqlPool_parseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mut  func(c *config.DatabaseConfig)
		want string
	}{
		{
			name: "invalid_sslmode_fails_pgx_parse",
			mut: func(c *config.DatabaseConfig) {
				c.SSLMode = "not_a_valid_ssl_mode"
			},
			want: "[pgsql/pool] failed to parse database config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t.Cleanup(cancel)

			cfg := testDatabaseConfig()
			tt.mut(cfg)
			_, err := NewPgsqlNative(ctx, cfg)
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

func TestPgsqlPool_HealthCheck_unreachable(t *testing.T) {
	t.Parallel()
	ctxPool, cancelPool := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancelPool)

	p, err := NewPgsqlNative(ctxPool, testDatabaseConfig())
	require.Nil(t, p)
	assertUnreachableDBError(t, err)
}
