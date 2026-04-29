package role

import (
	"ndinhbang/go-template/internal/domain/entity"
	"time"
)

type CreateRoleResponse struct {
	CreatedAt   time.Time            `json:"created_at"`
	Permissions entity.PermissionMap `json:"permissions"`
	Slug        string               `json:"slug"`
	Name        string               `json:"name"`
	ID          int64                `json:"id"`
	Lv          int                  `json:"lv"`
}
