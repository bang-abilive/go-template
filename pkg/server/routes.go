package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s *Server) SetupRoutes(routes []echo.Route) {
	s.echo.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the API!")
	})
	s.echo.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	for _, route := range routes {
		s.echo.AddRoute(route)
	}
}
