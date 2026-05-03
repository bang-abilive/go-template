package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s *Server) SetupRoutes(rr RouteRegistrar) {
	s.echo.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the API!")
	})
	s.echo.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	s.echo.GET("/debug/routes", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, s.echo.Router().Routes())
	})

	if rr != nil {
		rr.RegisterRoutes(s.echo)
	}
}
