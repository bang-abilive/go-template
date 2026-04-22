package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s *Server) SetupRoutes() {
	s.echo.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the API!")
	})
	s.echo.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	s.echo.GET("/api/v1/hello", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
