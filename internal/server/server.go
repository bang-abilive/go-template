package server

import (
	"log/slog"
	"net/http"

	"ndinhbang/go-skeleton/internal/config"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	MaxBodyLimitBytes int64 = 1024 * 1024 // 1MB
)

type Server struct {
	cfg  *config.ServerConfig
	echo *echo.Echo
}

func New(cfg *config.ServerConfig) *Server {
	return &Server{
		cfg: cfg,
		echo: echo.NewWithConfig(echo.Config{
			Logger: slog.Default(),
		}),
	}
}

func (s *Server) SetupMiddlewares(middlewares ...echo.MiddlewareFunc) {
	if len(middlewares) == 0 {
		// Default middlewares
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

func (s *Server) Start() error {
	// Setup routes
	if err := s.SetupRoutes(); err != nil {
		return err
	}

	slog.Info("starting server on address", "address", s.cfg.ServerAddress())

	return s.echo.Start(s.cfg.ServerAddress())
}
