package postgres

import (
	"context"
	"ndinhbang/go-template/internal/domain/entity"
	"ndinhbang/go-template/internal/usecase/role"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRoleRepository struct {
	db *pgxpool.Pool
}

func NewPgxRoleRepository(db *pgxpool.Pool) role.RoleRepository {
	return &pgxRoleRepository{db: db}
}

func (r *pgxRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (slug, name, lv, permissions) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, role.Slug, role.Name, role.Lv, role.Permissions).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
}

// Count implements [role.RoleRepository].
func (r *pgxRoleRepository) Count(ctx context.Context) (int64, error) {
	panic("unimplemented")
}

// Delete implements [role.RoleRepository].
func (r *pgxRoleRepository) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// Exists implements [role.RoleRepository].
func (r *pgxRoleRepository) Exists(ctx context.Context, id int64) (bool, error) {
	panic("unimplemented")
}

// Find implements [role.RoleRepository].
func (r *pgxRoleRepository) Find(ctx context.Context, id int64) (*entity.Role, error) {
	panic("unimplemented")
}

// Search implements [role.RoleRepository].
func (r *pgxRoleRepository) Search(ctx context.Context, query string) ([]entity.Role, error) {
	panic("unimplemented")
}

// Update implements [role.RoleRepository].
func (r *pgxRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	panic("unimplemented")
}
