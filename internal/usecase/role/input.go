package role

import "ndinhbang/go-template/internal/domain/entity"

type CreateRoleRequest struct {
	Permissions entity.PermissionMap `json:"permissions" validate:"required,dive,keys,alpha_dash,min=3,max=64,endkeys"`
	Slug        string               `json:"slug" validate:"omitempty,alpha_dash,min=3,max=64"`
	Name        string               `json:"name" validate:"required,min=3,max=64"`
	Lv          int                  `json:"lv" validate:"required,min=0,max=100"`
}
