package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ndinhbang/go-skeleton/internal/app"
	"ndinhbang/go-skeleton/pkg/config"
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
	return runContext(ctx)
}

// runContext loads config, builds the app, and runs until ctx is done.
// Called from [run] and from tests in this package.
func runContext(ctx context.Context) error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}
	application := app.NewApp(cfg)
	return application.Run(ctx)
}
