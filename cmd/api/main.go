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
	// Load the configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}
	// Setup signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// Stop listening for signals when the context is done
	defer stop()
	// Create the application
	app := app.NewApp(cfg)
	// Run the application
	return app.Run(ctx)
}
