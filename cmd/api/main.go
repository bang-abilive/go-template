package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ndinhbang/go-template/database/seeders"
	"ndinhbang/go-template/internal/app"
)

func main() {
	if err := run(); err != nil {
		slog.Error("[main] failed to start", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// parse flags
	seedFlag := flag.Bool("seed", false, "seed the database")
	flag.Parse()
	if *seedFlag {
		return runSeeder(ctx)
	}
	return runApi(ctx)
}

// runApi builds the app via DI and runs until ctx is done.
// Called from [run] and from tests in this package.
func runApi(ctx context.Context) error {
	application, err := app.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("[main] initialize app: %w", err)
	}
	return application.Run(ctx)
}

// runSeeder initializes the seeder and runs it.
func runSeeder(ctx context.Context) error {
	seeder, err := seeders.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("[main] initialize seeder: %w", err)
	}
	if err := seeder.Run(ctx); err != nil {
		return fmt.Errorf("[main] seed the database: %w", err)
	}
	slog.Info("[main] database seeded successfully")
	return nil
}
