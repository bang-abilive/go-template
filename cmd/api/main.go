package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"ndinhbang/go-skeleton/internal/config"
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

	// print go env
	slog.Info("[go] env", "env", os.Environ())

	// Print all config
	slog.Info("[config] app name", "name", cfg.App.Name)
	slog.Info("[config] app version", "version", cfg.App.Version)
	slog.Info("[config] app env", "env", cfg.App.Env)
	slog.Info("[config] app debug", "debug", cfg.App.Debug)
	slog.Info("[config] server address", "address", cfg.ServerAddress())
	slog.Info("[config] database name", "name", cfg.Database.Name)
	slog.Info("[config] database host", "host", cfg.Database.Host)
	slog.Info("[config] database port", "port", cfg.Database.Port)
	slog.Info("[config] database user", "user", cfg.Database.User)
	slog.Info("[config] database password", "password", cfg.Database.Password)
	slog.Info("[config] database ssl mode", "ssl mode", cfg.Database.SSLMode)
	slog.Info("[config] database max conns", "max conns", cfg.Database.MaxConns)
	slog.Info("[config] database max idle conns", "max idle conns", cfg.Database.MaxIdleConns)
	slog.Info("[config] database max lifetime conns", "max lifetime conns", cfg.Database.MaxLifetimeConns)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	// e.Use(middleware.Recover())       // recover panics as errors for proper error handling

	// Routes
	e.GET("/", hello)

	// Start server
	if err := e.Start(cfg.ServerAddress()); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
