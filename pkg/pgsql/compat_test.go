package pgsql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// *pgsqlCompat must implement [database.Database] (see ndinhbang/go-template/internal/database).
// var _ database.Database = (*pgsqlCompat)(nil)

func TestNewPgsqlCompat_HealthCheck_unreachable(t *testing.T) {
	t.Parallel()
	p, err := NewPgsqlCompat(context.Background(), testDatabaseConfig())
	require.NoError(t, err)
	require.NotNil(t, p)
	t.Cleanup(func() { _ = p.Close(context.Background()) })

	pingCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	t.Cleanup(cancel)

	_, err = p.Health(pingCtx)
	assertUnreachableDBError(t, err)
}

func TestPgsqlCompat_Close_idempotent(t *testing.T) {
	t.Parallel()
	p, err := NewPgsqlCompat(context.Background(), testDatabaseConfig())
	require.NoError(t, err)
	require.NoError(t, p.Close(context.Background()))
	require.NoError(t, p.Close(context.Background()))
}
