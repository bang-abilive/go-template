package handlers

import (
	"ndinhbang/go-template/internal/usecase/role"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

var _ server.GroupRegistrar = (*RoleHandler)(nil)

type RoleHandler struct {
	service role.Service
}

func NewRoleHandler(svc role.Service) *RoleHandler {
	return &RoleHandler{service: svc}
}

func (r *RoleHandler) RegisterRoutes(e *echo.Group) {
	g := e.Group("/role")

	g.AddRoute(echo.Route{Method: "POST", Path: "/create", Handler: r.Create, Name: "role.create"})
}

func (r *RoleHandler) Create(ctx *echo.Context) error {
	return nil
}
