package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	AppEnv     string `koanf:"APP_ENV"`
	AppName    string `koanf:"APP_NAME"`
	AppVersion string `koanf:"APP_VERSION"`

	ServerPort string `koanf:"SERVER_PORT"`

	DatabaseHost             string        `koanf:"DB_HOST"`
	DatabasePort             string        `koanf:"DB_PORT"`
	DatabaseUser             string        `koanf:"DB_USER"`
	DatabasePassword         string        `koanf:"DB_PASSWORD"`
	DatabaseName             string        `koanf:"DB_NAME"`
	DatabaseSSLMode          string        `koanf:"DB_SSL_MODE"`
	DatabaseMaxConns         int           `koanf:"DB_MAX_CONNS"`
	DatabaseMaxIdleConns     int           `koanf:"DB_MAX_IDLE_CONNS"`
	DatabaseMaxLifetimeConns time.Duration `koanf:"DB_MAX_LIFETIME_CONNS"`
}

const (
	defaultServerPort = "8080"
	envDevelopment    = "development"
	envTest           = "test"
	envStaging        = "staging"
	envProduction     = "production"
	defaultAppEnv     = envDevelopment
)

var validAppEnvs = []string{
	envDevelopment,
	envTest,
	envStaging,
	envProduction,
}

var allowedAppEnvs = map[string]struct{}{
	envDevelopment: {},
	envTest:        {},
	envStaging:     {},
	envProduction:  {},
}

func Load() (*Config, error) {
	mode := strings.TrimSpace(os.Getenv("APP_ENV"))
	if mode == "" {
		mode = defaultAppEnv
	}

	k := koanf.New(".")
	loadOrder := []string{
		".env",
		".env.local",
		fmt.Sprintf(".env.%s", mode),
		fmt.Sprintf(".env.%s.local", mode),
	}

	for _, path := range loadOrder {
		if err := k.Load(file.Provider(path), dotenv.Parser()); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("load env file %q: %w", path, err)
		}
	}

	cfg := &Config{
		AppEnv:     mode,
		ServerPort: defaultServerPort,
	}

	if err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{
		Tag:       "koanf",
		FlatPaths: true,
	}); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.AppEnv = strings.TrimSpace(cfg.AppEnv)
	if cfg.AppEnv == "" {
		cfg.AppEnv = mode
	}
	if _, ok := allowedAppEnvs[cfg.AppEnv]; !ok {
		return nil, fmt.Errorf("invalid APP_ENV %q: must be one of %s", cfg.AppEnv, strings.Join(validAppEnvs, "|"))
	}

	cfg.ServerPort = strings.TrimSpace(cfg.ServerPort)
	if cfg.ServerPort == "" {
		cfg.ServerPort = defaultServerPort
	}
	port, err := strconv.Atoi(cfg.ServerPort)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT %q: must be an integer", cfg.ServerPort)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid SERVER_PORT %q: must be between 1 and 65535", cfg.ServerPort)
	}

	return cfg, nil
}

func (c Config) ServerAddress() string {
	return ":" + c.ServerPort
}

func (c Config) IsProduction() bool {
	return c.AppEnv == envProduction
}

func (c Config) IsDevelopment() bool {
	return c.AppEnv == envDevelopment
}

func (c Config) IsTest() bool {
	return c.AppEnv == envTest
}

func (c Config) IsStaging() bool {
	return c.AppEnv == envStaging
}
