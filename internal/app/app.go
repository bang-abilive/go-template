package app

import (
	"context"
	"errors"

	"net/http"

	v1 "ndinhbang/go-template/internal/delivery/http/v1"
	"ndinhbang/go-template/internal/repository/postgres"
	"ndinhbang/go-template/internal/usecase/user"
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/pgsql"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

type App struct {
	cfg *config.Config
	srv *server.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run(ctx context.Context) error {
	db, err := pgsql.NewPgsqlNative(ctx, &a.cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	userRepo := postgres.NewPgxUserRepository(db.Pool())
	userService := user.NewService(userRepo)
	userHandler := v1.NewUserHandler(userService)

	a.srv = server.New(&a.cfg.Server)
	a.srv.SetupMiddlewares()

	routes := []echo.Route{
		{
			Method:  "POST",
			Path:    "/api/v1/users/register",
			Handler: userHandler.Register,
		},
	}

	a.srv.SetupRoutes(routes)

	if err := a.srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
