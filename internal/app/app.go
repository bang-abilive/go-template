package app

import (
	"context"
	"errors"

	"net/http"

	"ndinhbang/go-skeleton/pkg/config"
	"ndinhbang/go-skeleton/pkg/server"
)

type App struct {
	cfg *config.Config
	srv *server.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run(ctx context.Context) error {
	a.srv = server.New(&a.cfg.Server)
	a.srv.SetupMiddlewares()
	a.srv.SetupRoutes()
	if err := a.srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
