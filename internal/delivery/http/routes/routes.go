package routes

import (
	v1 "ndinhbang/go-template/internal/delivery/http/v1"
	v2 "ndinhbang/go-template/internal/delivery/http/v2"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

var _ server.RouteRegistrar = (*Registrar)(nil)

type Registrar struct {
	groups []server.GroupRegistrar
}

func NewRegistrar(
	v1rr *v1.RouteRegistrar,
	v2rr *v2.RouteRegistrar,
) *Registrar {
	return &Registrar{
		groups: []server.GroupRegistrar{
			v1rr,
			v2rr,
		},
	}
}

func (rr *Registrar) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/api")
	for _, group := range rr.groups {
		group.RegisterRoutes(g)
	}
}
