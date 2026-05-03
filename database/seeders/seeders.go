package seeders

import (
	"context"
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/db"
)

type Seeder struct {
	cfg *config.DatabaseConfig
	db  *db.PostgresDatabase
}

func NewSeeder(cfg *config.DatabaseConfig, db *db.PostgresDatabase) *Seeder {
	return &Seeder{cfg: cfg, db: db}
}

func (s *Seeder) Run(ctx context.Context) error {
	if err := s.seedUsers(ctx); err != nil {
		return err
	}
	if err := s.seedRoles(ctx); err != nil {
		return err
	}
	if err := s.seedPolicies(ctx); err != nil {
		return err
	}
	return nil
}
