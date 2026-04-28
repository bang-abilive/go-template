package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	return runContext(ctx)
}

// runContext builds the app via DI and runs until ctx is done.
// Called from [run] and from tests in this package.
func runContext(ctx context.Context) error {
	application, err := app.Initialize(ctx)
	if err != nil {
		return err
	}
	return application.Run(ctx)
}
