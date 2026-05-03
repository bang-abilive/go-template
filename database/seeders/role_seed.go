package seeders

import (
	"context"
	"fmt"
	"log/slog"
)

type seedRoleRow struct {
	slug        string
	name        string
	lv          int
	userEmail   string
	permissions string
}

// seedRoles inserts seed roles and user-role assignments.
// Roles are upserted by slug. user_role rows are inserted with ON CONFLICT DO NOTHING.
func (s *Seeder) seedRoles(ctx context.Context) error {
	roles := []seedRoleRow{
		{slug: "admin", name: "Administrator", lv: 100, userEmail: "admin@example.com", permissions: `{"admin": true}`},
		{slug: "editor", name: "Editor", lv: 50, userEmail: "editor@example.com", permissions: `{"editor": true}`},
		{slug: "viewer", name: "Viewer", lv: 10, userEmail: "viewer@example.com", permissions: `{"guest": true}`},
	}

	for _, r := range roles {
		// Upsert role by slug.
		var roleID int64
		err := s.db.Pool().QueryRow(ctx,
			`INSERT INTO roles (slug, name, lv, permissions)
			 VALUES ($1, $2, $3, $4::jsonb)
			 ON CONFLICT (slug) DO UPDATE
			   SET name = EXCLUDED.name, lv = EXCLUDED.lv, permissions = EXCLUDED.permissions
			 RETURNING id`,
			r.slug, r.name, r.lv, r.permissions,
		).Scan(&roleID)
		if err != nil {
			return fmt.Errorf("[seeders] seed role %s: %w", r.slug, err)
		}
		slog.Info("[seeders] role seeded", "slug", r.slug, "id", roleID)

		// Look up the user id.
		var userID int64
		err = s.db.Pool().QueryRow(ctx,
			`SELECT id FROM users WHERE email = $1`,
			r.userEmail,
		).Scan(&userID)
		if err != nil {
			return fmt.Errorf("[seeders] find user %s: %w", r.userEmail, err)
		}

		// Assign user → role (idempotent).
		_, err = s.db.Pool().Exec(ctx,
			`INSERT INTO user_role (user_id, role_id) VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`,
			userID, roleID,
		)
		if err != nil {
			return fmt.Errorf("[seeders] assign user %d to role %s: %w", userID, r.slug, err)
		}
		slog.Info("[seeders] user-role assigned", "user_id", userID, "role", r.slug)
	}
	return nil
}
