package server

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v5"

	"ndinhbang/go-skeleton/pkg/config"
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
