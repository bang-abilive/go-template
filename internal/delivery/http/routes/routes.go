package routes

import (
	"ndinhbang/go-template/internal/delivery/http/middleware"
	v1 "ndinhbang/go-template/internal/delivery/http/v1"
	"ndinhbang/go-template/internal/delivery/http/v1/handlers"
	v2 "ndinhbang/go-template/internal/delivery/http/v2"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

var _ server.RouteRegistrar = (*Registrar)(nil)

type Registrar struct {
	v1rr         *v1.RouteRegistrar
	v2rr         *v2.RouteRegistrar
	authMW       *middleware.AuthMiddleware
	casbinMW     *middleware.CasbinMiddleware
	authorizeHdl *handlers.AuthorizeHandler
}

func NewRegistrar(
	v1rr *v1.RouteRegistrar,
	v2rr *v2.RouteRegistrar,
	authMW *middleware.AuthMiddleware,
	casbinMW *middleware.CasbinMiddleware,
	authorizeHdl *handlers.AuthorizeHandler,
) *Registrar {
	return &Registrar{
		v1rr:         v1rr,
		v2rr:         v2rr,
		authMW:       authMW,
		casbinMW:     casbinMW,
		authorizeHdl: authorizeHdl,
	}
}

func (rr *Registrar) RegisterRoutes(e *echo.Echo) {
	// Public group: no auth required. Used for the debug authorize endpoint.
	public := e.Group("/api/v1")
	rr.authorizeHdl.RegisterRoutes(public)

	// Protected group: every request must carry ?user_id and pass ABAC check.
	protected := e.Group("/api",
		rr.authMW.Middleware(),
		rr.casbinMW.Middleware(),
	)
	rr.v1rr.RegisterRoutes(protected)
	rr.v2rr.RegisterRoutes(protected)
}
