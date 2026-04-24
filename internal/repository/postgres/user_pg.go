package postgres

import (
	"context"
	"errors"
	"ndinhbang/go-skeleton/internal/domain/entity"
	"ndinhbang/go-skeleton/internal/usecase/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository(db *pgxpool.Pool) user.UserRepository {
	return &pgxUserRepository{db: db}
}

func (r *pgxUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, user.Email.Value(), user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *pgxUserRepository) Find(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	var u entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgxUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	var u entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
