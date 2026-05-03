package seeders

import (
	"context"
	"fmt"
	"log/slog"
)

type seedUserRow struct {
	email    string
	password string
}

// seedUsers inserts seed users into the users table.
// Idempotent: existing emails are updated in-place (DO UPDATE) so RETURNING always yields the id.
func (s *Seeder) seedUsers(ctx context.Context) error {
	users := []seedUserRow{
		{email: "admin@example.com", password: "password"},
		{email: "editor@example.com", password: "password"},
		{email: "viewer@example.com", password: "password"},
	}

	for _, u := range users {
		var id int64
		err := s.db.Pool().QueryRow(ctx,
			`INSERT INTO users (email, password)
			 VALUES ($1, $2)
			 ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
			 RETURNING id`,
			u.email, u.password,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("[seeders] seed user %s: %w", u.email, err)
		}
		slog.Info("[seeders] user seeded", "email", u.email, "id", id)
	}
	return nil
}
