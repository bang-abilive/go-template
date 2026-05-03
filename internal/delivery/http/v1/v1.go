package v1

import (
	"ndinhbang/go-template/internal/delivery/http/v1/handlers"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

var _ server.GroupRegistrar = (*RouteRegistrar)(nil)

// RouteRegistrar is responsible for registering all v1 routes.
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
	v1 := g.Group("/v1")
	for _, rr := range h.groups {
		rr.RegisterRoutes(v1)
	}
}
