package authorizer

import (
	"fmt"
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/db"

	"github.com/casbin/casbin/v3"
	pgxadapter "github.com/noho-digital/casbin-pgx-adapter"
)

type Authorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(db *db.PostgresDatabase, dbName string, tableName string) (*Authorizer, error) {
	// Create the adapter with optional configuration
	adapter, err := pgxadapter.NewAdapterWithPool(db.GetPool(),
		pgxadapter.WithDatabaseName(dbName),
		pgxadapter.WithTableName(tableName),
		// pgxadapter.WithIndex("ptype", "v0"),
	)

	if err != nil {
		return nil, fmt.Errorf("[authorizer] create adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer("pkg/authorizer/model/rbac.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("[authorizer] create enforcer: %w", err)
	}

	// Load the policy from DB.
	if err = enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("[authorizer] load policy: %w", err)
	}

	return &Authorizer{
		enforcer: enforcer,
	}, nil
}

func NewDefaultAuthorizer(db *db.PostgresDatabase, cfg *config.DatabaseConfig) (*Authorizer, error) {
	return NewAuthorizer(db, cfg.Name, "policies")
}

func (a *Authorizer) GetEnforcer() *casbin.Enforcer {
	return a.enforcer
}
