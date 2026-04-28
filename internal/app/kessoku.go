//go:generate go tool kessoku $GOFILE
package app

import (
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
	kessoku.Provide(NewApp),
)