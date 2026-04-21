package main

import (
	"log/slog"
	"os"
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

	srv := server.New(&cfg.Server)
	srv.SetupMiddlewares()

	if err := srv.Start(); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
