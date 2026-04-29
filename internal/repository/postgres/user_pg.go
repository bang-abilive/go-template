package postgres

import (
	"context"
	"errors"
	"fmt"
	"ndinhbang/go-template/internal/domain/entity"
	"ndinhbang/go-template/internal/usecase/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository(db *pgxpool.Pool) user.UserRepository {
	return &pgxUserRepository{db: db}
}

func (r *pgxUserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("[repository] count user: %w", err)
	}
	return count, nil
}

func (r *pgxUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("[repository] delete user (id=%d): %w", id, err)
	}
	return nil
}

func (r *pgxUserRepository) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("[repository] exists user (id=%d): %w", id, err)
	}
	return exists, nil
}

// FindByEmailAndPassword implements [user.UserRepository].
func (r *pgxUserRepository) FindByEmailAndPassword(ctx context.Context, email string, password string) (*entity.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1 AND password = $2`
	var u entity.User
	err := r.db.QueryRow(ctx, query, email, password).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[repository] find user by email and password (email=%s): %w", email, err)
	}
	return &u, nil
}

// Search implements [user.UserRepository].
func (r *pgxUserRepository) Search(ctx context.Context, query string) ([]entity.User, error) {
	panic("unimplemented")
}

func (r *pgxUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, user.Email.Value(), user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *pgxUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET email = $1, password = $2, updated_at = NOW() WHERE id = $3 RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, user.Email.Value(), user.Password, user.ID).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
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
