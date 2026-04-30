// Package entity contains the domain entities for the application.
package entity

import (
	"ndinhbang/go-template/internal/domain/values"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     values.Email
	Password  string
	ID        int64
}
