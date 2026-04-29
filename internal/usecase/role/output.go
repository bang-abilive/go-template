package role

import (
	"ndinhbang/go-template/internal/domain/entity"
	"time"
)

type CreateRoleResponse struct {
	ID          int64              `json:"id"`
	Slug        string             `json:"slug"`
	Name        string             `json:"name"`
	Lv          int                `json:"lv"`
	Permissions entity.Permissions `json:"permissions"`
	CreatedAt   time.Time          `json:"created_at"`
}
