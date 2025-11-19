package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer string `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" env-default:"http://localhost:8080"`
	DB_DSN     string `env:"DATABASE_DSN" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		// TODO: add ReadEnv
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}
