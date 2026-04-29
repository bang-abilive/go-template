package role

import (
	"context"
	"ndinhbang/go-template/internal/domain/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, id int64) (*entity.Role, error)
	Search(ctx context.Context, query string) ([]entity.Role, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context, id int64) (bool, error)
}
