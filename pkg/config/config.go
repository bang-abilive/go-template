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
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION" envDefault:"0.0.1"`
	Debug   bool   `env:"APP_DEBUG" envDefault:"false"`
}

type DatabaseConfig struct {
	Name            string        `env:"DB_NAME,notEmpty"`
	Schema          string        `env:"DB_SCHEMA,notEmpty" envDefault:"public"`
	Host            string        `env:"DB_HOST,notEmpty"`
	User            string        `env:"DB_USER,notEmpty,unset"`
	Password        string        `env:"DB_PASSWORD,notEmpty,unset"`
	SSLMode         string        `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"10m"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"20"`
	MinConns        int32         `env:"DB_MIN_CONNS" envDefault:"5"`
	MaxIdleConns    int32         `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	Port            uint16        `env:"DB_PORT,notEmpty" envDefault:"5432"`
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
			slog.Warn("[config] not found:", "file", file)
			continue
		}
		if err := godotenv.Load(file); err != nil {
			slog.Warn("[config] failed to load:", "file", file, "error", err)
			continue
		}
		slog.Info("[config] loaded:", "file", file)
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("[config] parse env: %w", err)
	}

	slog.Info("[config] environment", "env", cfg.App.Env)

	return &cfg, nil
}

func GetServerConfig(cfg *Config) *ServerConfig {
	return &cfg.Server
}

func (s ServerConfig) ServerAddress() string {
	return fmt.Sprintf(":%d", s.Port)
}

func GetDatabaseConfig(cfg *Config) *DatabaseConfig {
	return &cfg.Database
}

func (c DatabaseConfig) DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s", c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode, c.Schema)
}
