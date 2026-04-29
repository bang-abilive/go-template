package entity

import (
	"ndinhbang/go-template/internal/domain/values"
	"time"
)

type PermissionMap map[string]bool

type Role struct {
	CreatedAt   time.Time
	Permissions PermissionMap
	Slug        values.Slug
	UpdatedAt   time.Time
	Name        string
	ID          int64
	Lv          int
}

type UserRole struct {
	UserID int64
	RoleID int64
}
