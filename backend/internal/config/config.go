package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env              string
	Port             string
	LogLevel         string
	GinMode          string
	PrimBaseURL      string
	PrimAPIKey       string
	APIKey           string
	SentryEnabled    bool
	SentryDSN        string
	SentryEnableLogs bool
}

// Load loads and validates the configuration.
// Returns an error when a required environment variable is missing.
func Load() (Config, error) {
	// Load .env from working directory if present
	_ = godotenv.Load()

	primBase, err := requireEnv("PRIM_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	primKey, err := requireEnv("PRIM_API_KEY")
	if err != nil {
		return Config{}, err
	}
	apiKey, err := requireEnv("API_KEY")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Env:              getEnv("ENV", "dev"),
		Port:             getEnv("PORT", "8080"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		GinMode:          getEnv("GIN_MODE", "debug"),
		PrimBaseURL:      primBase,
		PrimAPIKey:       primKey,
		APIKey:           apiKey,
		SentryEnabled:    getEnv("SENTRY_ENABLED", "false") == "true",
		SentryDSN:        getEnv("SENTRY_DSN", ""),
		SentryEnableLogs: getEnv("SENTRY_ENABLE_LOGS", "false") == "true",
	}
	return cfg, nil
}

// getEnv returns the value of the environment variable or a default
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// requireEnv returns the value of the environment variable or an error if missing
func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required env variable: %s", key)
	}
	return val, nil
}
