package entity

import (
	"ndinhbang/go-skeleton/internal/domain/values"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     values.Email
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
