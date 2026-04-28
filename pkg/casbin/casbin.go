package casbin

import (
	"fmt"
	"ndinhbang/go-template/pkg/config"

	"github.com/casbin/casbin/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxadapter "github.com/noho-digital/casbin-pgx-adapter"
)

type Gatekeeper struct {
	enforcer *casbin.Enforcer
}

func NewGatekeeper(cfg *config.DatabaseConfig, dbpool *pgxpool.Pool, modelPath string) (*Gatekeeper, error) {
	// Create the adapter with optional configuration
	adapter, err := pgxadapter.NewAdapterWithPool(dbpool,
		pgxadapter.WithTableName("policies"),  // Optional: custom table name
		pgxadapter.WithDatabaseName(cfg.Name), // Optional: custom database name
	)

	if err != nil {
		return nil, fmt.Errorf("[casbin] failed to create adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("[casbin] failed to create enforcer: %w", err)
	}

	// // Load the policy from DB.
	// if err = enforcer.LoadPolicy(); err != nil {
	// 	return nil, fmt.Errorf("[casbin] failed to load policy: %w", err)
	// }

	return &Gatekeeper{
		enforcer: enforcer,
	}, nil
}

func (g *Gatekeeper) LoadPolicy() error {
	return g.enforcer.LoadPolicy()
}

func (g *Gatekeeper) SavePolicy() error {
	return g.enforcer.SavePolicy()
}

func (g *Gatekeeper) Enforce(sub, obj, act string) (bool, error) {
	return g.enforcer.Enforce(sub, obj, act)
}

func (g *Gatekeeper) EnforceEx(sub, obj, act string) (bool, []string, error) {
	return g.enforcer.EnforceEx(sub, obj, act)
}

func (g *Gatekeeper) AddPolicy(sub, obj, act string) (bool, error) {
	return g.enforcer.AddPolicy(sub, obj, act)
}

func (g *Gatekeeper) CreatePolicy(sub, obj, act string) (bool, error) {
	return g.enforcer.AddPolicy(sub, obj, act)
}

func (g *Gatekeeper) AddGroupingPolicy(sub, role string) (bool, error) {
	return g.enforcer.AddGroupingPolicy(sub, role)
}

func (g *Gatekeeper) RemovePolicy(sub, obj, act string) (bool, error) {
	return g.enforcer.RemovePolicy(sub, obj, act)
}

func (g *Gatekeeper) RemoveGroupingPolicy(sub, obj, act string) (bool, error) {
	return g.enforcer.RemoveGroupingPolicy(sub, obj, act)
}
