package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetAndLoadAPIKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("BASELINE_CONFIG", path)
	t.Setenv("BASELINE_API_KEY", "")

	if err := SetAPIKey("api_test"); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "api_test" {
		t.Fatalf("expected api_test, got %q", cfg.APIKey)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestAPIKeyPrefersEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("BASELINE_CONFIG", path)
	t.Setenv("BASELINE_API_KEY", "")

	if err := SetAPIKey("api_config"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BASELINE_API_KEY", "api_env")

	value, source, err := APIKey()
	if err != nil {
		t.Fatal(err)
	}
	if value != "api_env" {
		t.Fatalf("expected api_env, got %q", value)
	}
	if source != "BASELINE_API_KEY" {
		t.Fatalf("expected env source, got %q", source)
	}
}

func TestMask(t *testing.T) {
	if got := Mask("api_1234567890"); got != "api_******7890" {
		t.Fatalf("unexpected mask: %q", got)
	}
}
