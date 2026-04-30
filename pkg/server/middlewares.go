// Package server contains the HTTP server implementation.
package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	MaxBodyLimitBytes int64 = 1024 * 1024 // 1MB
)

func (s *Server) SetupMiddlewares(middlewares ...echo.MiddlewareFunc) {
	if len(middlewares) == 0 {
		middlewares = []echo.MiddlewareFunc{
			middleware.RequestLogger(),
			middleware.Secure(),
			middleware.Recover(),
			middleware.BodyLimit(MaxBodyLimitBytes),
			middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
				AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			}),
		}
	}
	for _, middleware := range middlewares {
		s.echo.Use(middleware)
	}
}
