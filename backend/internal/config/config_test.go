package config

import (
	"os"
	"testing"
)

func TestLoadFrom_ValidConfig(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(`
env: test
port: "9090"
log_level: warn
gin_mode: release
prim:
  base_url: https://example.com
  api_key: primkey123
api_key: myapikey456
sentry:
  dsn: https://abc@sentry.io/1
  enable_logs: true
`)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := LoadFrom(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Env != "test" {
		t.Errorf("expected env=test, got %q", cfg.Env)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected port=9090, got %q", cfg.Port)
	}
	if cfg.Prim.BaseURL != "https://example.com" {
		t.Errorf("expected prim.base_url=https://example.com, got %q", cfg.Prim.BaseURL)
	}
	if cfg.Prim.APIKey != "primkey123" {
		t.Errorf("expected prim.api_key=primkey123, got %q", cfg.Prim.APIKey)
	}
	if cfg.APIKey != "myapikey456" {
		t.Errorf("expected api_key=myapikey456, got %q", cfg.APIKey)
	}
	if cfg.Sentry.DSN != "https://abc@sentry.io/1" {
		t.Errorf("expected sentry.dsn set, got %q", cfg.Sentry.DSN)
	}
	if !cfg.Sentry.EnableLogs {
		t.Error("expected sentry.enable_logs=true")
	}
}

func TestLoadFrom_MissingPrimBaseURL(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(`
prim:
  api_key: primkey123
api_key: myapikey456
`)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	_, err = LoadFrom(f.Name())
	if err == nil {
		t.Fatal("expected error for missing prim.base_url, got nil")
	}
}

func TestLoadFrom_MissingPrimAPIKey(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(`
prim:
  base_url: https://example.com
api_key: myapikey456
`)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	_, err = LoadFrom(f.Name())
	if err == nil {
		t.Fatal("expected error for missing prim.api_key, got nil")
	}
}

func TestLoadFrom_MissingAPIKey(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(`
prim:
  base_url: https://example.com
  api_key: primkey123
`)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	_, err = LoadFrom(f.Name())
	if err == nil {
		t.Fatal("expected error for missing api_key, got nil")
	}
}

func TestLoadFrom_Defaults(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(`
prim:
  base_url: https://example.com
  api_key: primkey123
api_key: myapikey456
`)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := LoadFrom(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Env != "dev" {
		t.Errorf("expected default env=dev, got %q", cfg.Env)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default port=8080, got %q", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log_level=info, got %q", cfg.LogLevel)
	}
	if cfg.GinMode != "debug" {
		t.Errorf("expected default gin_mode=debug, got %q", cfg.GinMode)
	}
}

func TestLoadFrom_FileNotFound(t *testing.T) {
	_, err := LoadFrom("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
