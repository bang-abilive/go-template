package pgsql

import (
	"context"
	"testing"
	"time"

	"ndinhbang/go-skeleton/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// *pgsqlConnection must implement [database.Database] (see ndinhbang/go-skeleton/internal/database).
var _ database.Database = (*pgsqlConnection)(nil)

func TestNewPgsqlConnection_refused(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	_, err := NewPgsqlConnection(ctx, testDatabaseConfig())
	require.Error(t, err)
	assert.ErrorContains(t, err, "[pgsql/db] failed to connect to database")
}
