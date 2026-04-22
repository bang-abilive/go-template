package server

import (
	"context"
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

// Serve starts the HTTP server with the given StartConfig.
// Use this directly when you need TLS, a custom Listener, or ListenerAddrFunc.
func (s *Server) Serve(ctx context.Context, sc echo.StartConfig) error {
	slog.Info("starting server", "address", sc.Address)
	return sc.Start(ctx, s.echo)
}

// Start serves with the default address derived from ServerConfig.
func (s *Server) Start(ctx context.Context) error {
	return s.Serve(ctx, echo.StartConfig{
		Address:    s.cfg.ServerAddress(),
		HideBanner: true,
		HidePort:   true,
	})
}
