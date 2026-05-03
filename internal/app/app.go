package app

import (
	"context"
	"errors"
	"ndinhbang/go-template/pkg/authorizer"
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/db"
	"ndinhbang/go-template/pkg/server"
	"net/http"
)

type App struct {
	cfg        *config.Config
	srv        *server.Server
	db         *db.PostgresDatabase
	authorizer *authorizer.Authorizer
	rr         server.RouteRegistrar
}

func NewApp(
	cfg *config.Config,
	srv *server.Server,
	database *db.PostgresDatabase,
	authorizer *authorizer.Authorizer,
	rr server.RouteRegistrar,
) *App {
	return &App{
		cfg:        cfg,
		srv:        srv,
		db:         database,
		authorizer: authorizer,
		rr:         rr,
	}
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()

	a.srv.SetupMiddlewares()
	a.srv.SetupRoutes(a.rr)

	if err := a.srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
