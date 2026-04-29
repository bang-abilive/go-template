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
	// Seed the database with the default data
	if err := s.seedUsers(ctx); err != nil {
		return err
	}
	return nil
}
