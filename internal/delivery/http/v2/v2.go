package v2

import (
	"ndinhbang/go-template/internal/delivery/http/v1/handlers"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

var _ server.GroupRegistrar = (*RouteRegistrar)(nil)

type RouteRegistrar struct {
	groups []server.GroupRegistrar
}

func New(
	userHdl *handlers.UserHandler,
	roleHdl *handlers.RoleHandler,
) *RouteRegistrar {
	return &RouteRegistrar{
		groups: []server.GroupRegistrar{
			userHdl,
			roleHdl,
		},
	}
}

func (h *RouteRegistrar) RegisterRoutes(g *echo.Group) {
	v2 := g.Group("/v2")
	for _, group := range h.groups {
		group.RegisterRoutes(v2)
	}
}
