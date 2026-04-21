package config

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env     string `env:"APP_ENV,notEmpty" envDefault:"development"`
	Debug   bool   `env:"APP_DEBUG" envDefault:"false"`
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION" envDefault:"0.0.1"`
}

type DatabaseConfig struct {
	Name             string        `env:"DB_NAME,notEmpty"`
	Host             string        `env:"DB_HOST,notEmpty"`
	Port             string        `env:"DB_PORT,notEmpty"`
	User             string        `env:"DB_USER,notEmpty,unset"`
	Password         string        `env:"DB_PASSWORD,notEmpty,unset"`
	SSLMode          string        `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxConns         int           `env:"DB_MAX_CONNS" envDefault:"10"`
	MaxIdleConns     int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	MaxLifetimeConns time.Duration `env:"DB_MAX_LIFETIME_CONNS" envDefault:"1h"`
}

type ServerConfig struct {
	Port uint16 `env:"SERVER_PORT,notEmpty" envDefault:"8080"`
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

// safeEnvName rejects any APP_ENV value that contains path-traversal characters.
// Only alphanumeric, hyphens, and underscores are allowed (e.g. "development", "prod-us").
var safeEnvName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// LoadFromEnv loads layered .env files into the process environment, then parses into [Config].
func LoadFromEnv() (*Config, error) {
	mode := strings.TrimSpace(os.Getenv("APP_ENV"))
	if mode == "" {
		mode = "development"
	}
	if !safeEnvName.MatchString(mode) {
		return nil, fmt.Errorf("[config] invalid APP_ENV value %q: only alphanumeric, hyphens, and underscores are allowed", mode)
	}

	candidates := []string{
		".env." + mode + ".local",
		".env." + mode,
		".env.local",
		".env",
	}

	for _, file := range candidates {
		if _, err := os.Stat(file); err != nil {
			continue
		}
		if err := godotenv.Load(file); err != nil {
			slog.Warn("[config] failed to load env file", "file", file, "error", err)
			continue
		}
		slog.Info("[config] loaded env file", "file", file)
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("[config] parse env: %w", err)
	}

	slog.Info("[config] environment", "env", cfg.App.Env)

	return &cfg, nil
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

// func (c Config) DatabaseDSN() string {
// 	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode)
// }
