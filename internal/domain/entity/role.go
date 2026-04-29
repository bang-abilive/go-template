package entity

import (
	"ndinhbang/go-template/internal/domain/values"
	"time"
)

type Permissions map[string]bool

type Role struct {
	ID          int64
	Slug        values.Slug
	Name        string
	Lv          int
	Permissions Permissions
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRole struct {
	UserID int64
	RoleID int64
}
