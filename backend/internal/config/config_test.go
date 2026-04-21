package config

import (
	"testing"
)

func TestLoadFrom_FileNotFound(t *testing.T) {
	_, err := LoadFrom("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
