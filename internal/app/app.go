package app

import (
	"context"
	"errors"
	"ndinhbang/go-template/pkg/authorizer"
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/db"
	"ndinhbang/go-template/pkg/server"
	"net/http"

	"github.com/labstack/echo/v5"
)

type App struct {
	cfg         *config.Config
	srv         *server.Server
	db          *db.PostgresDatabase
	authorizer  *authorizer.Authorizer
	// Add fields for dependencies, e.g., database, server, etc.
}

func NewApp(
	cfg *config.Config,
	srv *server.Server,
	database *db.PostgresDatabase,
	authorizer *authorizer.Authorizer,
) *App {
	return &App{
		cfg:        cfg,
		srv:        srv,
		db:         database,
		authorizer: authorizer,
	}
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()
	// Start the server and other components here
	// a.authorizer.AddPolicy("alice", "data1", "read")
	// a.authorizer.AddPolicy("bob", "data2", "write")
	// a.authorizer.AddGroupingPolicy("alice", "admin")

	// if err := a.authorizer.SavePolicy(); err != nil {
	// 	return fmt.Errorf("[casbin] failed to save policy: %w", err)
	// }

	// if allowed, explanation, _ := a.gatekeeper.EnforceEx("bob", "data1", "read"); allowed {
	// 	log.Println("Alice can read data1")
	// } else {
	// 	log.Println("Bob cannot read data1")
	// 	log.Println(explanation)
	// }

	a.srv.SetupMiddlewares()
	a.srv.SetupRoutes([]echo.Route{
		{

		},
	})

	if err := a.srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}