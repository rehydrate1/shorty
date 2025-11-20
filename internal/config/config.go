package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer string `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" env-default:"http://localhost:8080"`
	DB_DSN     string `env:"DATABASE_DSN" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("failed to read env vars: %w", err)
		}
	} else {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	return &cfg, nil
}
