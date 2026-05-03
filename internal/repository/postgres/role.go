package postgres

import (
	"context"
	"ndinhbang/go-template/internal/domain/entity"
	"ndinhbang/go-template/pkg/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *db.PostgresDatabase) *RoleRepository {
	return &RoleRepository{db: db.Pool()}

}

func (r *RoleRepository) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (slug, name, lv, permissions) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, role.Slug, role.Name, role.Lv, role.Permissions).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
}

// Count implements [role.Repository].
func (r *RoleRepository) Count(ctx context.Context) (int64, error) {
	panic("unimplemented")
}

// Delete implements [role.Repository].
func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// Exists implements [role.Repository].
func (r *RoleRepository) Exists(ctx context.Context, id int64) (bool, error) {
	panic("unimplemented")
}

// Find implements [role.Repository].
func (r *RoleRepository) Find(ctx context.Context, id int64) (*entity.Role, error) {
	panic("unimplemented")
}

// Search implements [role.Repository].
func (r *RoleRepository) Search(ctx context.Context, query string) ([]entity.Role, error) {
	panic("unimplemented")
}

// Update implements [role.Repository].
func (r *RoleRepository) Update(ctx context.Context, role *entity.Role) error {
	panic("unimplemented")
}
