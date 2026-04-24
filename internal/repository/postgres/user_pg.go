package postgres

import (
	"context"
	"ndinhbang/go-skeleton/internal/domain/entity"
	"ndinhbang/go-skeleton/internal/usecase/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository(db *pgxpool.Pool) user.UserRepository {
	return &pgxUserRepository{db: db}
}

func (r *pgxUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, email, password, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password, user.CreatedAt)
	return err
}

func (r *pgxUserRepository) Find(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE id = $1`
	var user entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	return nil, err
}

func (r *pgxUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`

	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err // Thực tế nên handle sql.ErrNoRows để trả về nil, nil
	}
	return &user, nil
}
