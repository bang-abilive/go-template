package user

import (
	"context"
	"ndinhbang/go-template/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Search(ctx context.Context, query string) ([]entity.User, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context, id int64) (bool, error)
	FindByEmailAndPassword(ctx context.Context, email string, password string) (*entity.User, error)
}
