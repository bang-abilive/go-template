//go:generate go tool kessoku $GOFILE
package seeders

import (
	"ndinhbang/go-template/pkg/config"
	"ndinhbang/go-template/pkg/db"

	"github.com/mazrean/kessoku"
)

var _ = kessoku.Inject[*Seeder](
	"Initialize",
	kessoku.Provide(config.LoadFromEnv),
	kessoku.Provide(config.GetDatabaseConfig),
	kessoku.Async(kessoku.Provide(db.NewPostgresDatabase)),
	kessoku.Provide(NewSeeder),
)
