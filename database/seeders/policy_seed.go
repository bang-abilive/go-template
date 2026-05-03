package seeders

import (
	"context"
	"fmt"
	"log/slog"
)

type seedPolicyRow struct {
	ptype string
	v0    string // sub_rule (ABAC expression) or role (for g-type)
	v1    string // object  (e.g. route path)
	v2    string // action  (e.g. HTTP method)
}

// seedPolicies inserts ABAC p-type policies into the policies table.
// ON CONFLICT DO NOTHING keeps the operation idempotent.
func (s *Seeder) seedPolicies(ctx context.Context) error {
	policies := []seedPolicyRow{
		// Admin-only: create roles
		{ptype: "p", v0: `r.sub.Role == "admin"`, v1: "/api/v1/role/create", v2: "POST"},
		// Editor and above (lv >= 50): register new users
		{ptype: "p", v0: "r.sub.Level >= 50", v1: "/api/v1/user/register", v2: "POST"},
		// Everyone (lv >= 0): use the debug authorize endpoint itself
		{ptype: "p", v0: "r.sub.Level >= 0", v1: "/api/v1/authorize", v2: "POST"},
	}

	for _, p := range policies {
		_, err := s.db.Pool().Exec(ctx,
			`INSERT INTO policies (ptype, v0, v1, v2)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT DO NOTHING`,
			p.ptype, p.v0, p.v1, p.v2,
		)
		if err != nil {
			return fmt.Errorf("[seeders] seed policy (v0=%s v1=%s v2=%s): %w", p.v0, p.v1, p.v2, err)
		}
		slog.Info("[seeders] policy seeded", "ptype", p.ptype, "v0", p.v0, "v1", p.v1, "v2", p.v2)
	}
	return nil
}
