package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Env      string `yaml:"env"`
	Port     string `yaml:"port"`
	LogLevel string `yaml:"log_level"`
	GinMode  string `yaml:"gin_mode"`
	APIKey   string `yaml:"api_key"`
	Prim     struct {
		BaseURL string `yaml:"base_url"`
		APIKey  string `yaml:"api_key"`
	} `yaml:"prim"`
	Sentry struct {
		Enabled    bool   `yaml:"enabled"`
		DSN        string `yaml:"dsn"`
		EnableLogs bool   `yaml:"enable_logs"`
	} `yaml:"sentry"`
}

func Load() (Config, error) {
	return LoadFrom("config.yaml")
}

func LoadFrom(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}

	if cfg.Env == "" {
		cfg.Env = "dev"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.GinMode == "" {
		cfg.GinMode = "debug"
	}

	if cfg.Prim.BaseURL == "" {
		return Config{}, fmt.Errorf("missing required config key: prim.base_url")
	}
	if cfg.Prim.APIKey == "" {
		return Config{}, fmt.Errorf("missing required config key: prim.api_key")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing required config key: api_key")
	}

	return cfg, nil
}
