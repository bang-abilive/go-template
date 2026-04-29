package entity

import "time"

type Role struct {
	ID        int64
	Slug      string
	Name      string
	Lv        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole struct {
	UserID int64
	RoleID int64
}
