package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ndinhbang/go-skeleton/internal/config"
	"ndinhbang/go-skeleton/internal/server"
)

func main() {
	loadStart := time.Now()
	cfg, err := config.LoadFromEnv()
	loadDuration := time.Since(loadStart)
	if err != nil {
		slog.Error("failed to load config", "error", err, "duration", loadDuration)
		os.Exit(1)
	}
	slog.Info("[config] load from env", "duration", loadDuration)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(&cfg.Server)
	srv.SetupMiddlewares()
	srv.SetupRoutes()

	if err := srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
