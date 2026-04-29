package seeders

import (
	"context"
	"fmt"
	"log/slog"
	"ndinhbang/go-template/internal/domain/entity"
	"ndinhbang/go-template/internal/domain/values"
)

func (s *Seeder) seedUsers(ctx context.Context) error {
	emailVO, err := values.NewEmail("admin@example.com")
	if err != nil {
		return fmt.Errorf("[seeders] create email value object: %w", err)
	}
	password := "password"
	user := &entity.User{
		Email:    emailVO,
		Password: password,
	}

	if err := s.db.GetPool().QueryRow(ctx, "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", user.Email.Value(), user.Password).Scan(&user.ID); err != nil {
		return fmt.Errorf("[seeders] seed users: %w", err)
	} else {
		slog.Info("[seeders] user seeded:", "email", user.Email.Value())
	}
	return nil
}
