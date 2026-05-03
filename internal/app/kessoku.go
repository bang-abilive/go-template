//go:generate go tool kessoku $GOFILE
package app

import (
	"ndinhbang/go-template/internal/delivery/http/middleware"
	routes "ndinhbang/go-template/internal/delivery/http/routes"
	v1 "ndinhbang/go-template/internal/delivery/http/v1"
	"ndinhbang/go-template/internal/delivery/http/v1/handlers"
	v2 "ndinhbang/go-template/internal/delivery/http/v2"
	"ndinhbang/go-template/internal/repository/postgres"
	"ndinhbang/go-template/internal/usecase/role"
	"ndinhbang/go-template/internal/usecase/user"
	"ndinhbang/go-template/pkg/authorizer"
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/db"
	"ndinhbang/go-template/pkg/server"

	"github.com/mazrean/kessoku"
)

var _ = kessoku.Inject[*App](
	"Initialize",
	kessoku.Provide(config.LoadFromEnv),
	kessoku.Provide(config.GetDatabaseConfig),
	kessoku.Provide(config.GetServerConfig),
	kessoku.Async(kessoku.Provide(db.NewPostgresDatabase)),
	kessoku.Provide(authorizer.NewDefaultAuthorizer),
	kessoku.Provide(server.New),
	kessoku.Bind[user.Repository](kessoku.Provide(postgres.NewUserRepository)),
	kessoku.Provide(user.NewService),
	kessoku.Provide(handlers.NewUserHandler),
	kessoku.Bind[role.Repository](kessoku.Provide(postgres.NewRoleRepository)),
	kessoku.Provide(role.NewService),
	kessoku.Provide(handlers.NewRoleHandler),
	kessoku.Provide(middleware.NewAuthMiddleware),
	kessoku.Provide(middleware.NewCasbinMiddleware),
	kessoku.Provide(handlers.NewAuthorizeHandler),
	kessoku.Provide(v1.New),
	kessoku.Provide(v2.New),
	kessoku.Bind[server.RouteRegistrar](kessoku.Provide(routes.NewRegistrar)),
	kessoku.Provide(NewApp),
)
