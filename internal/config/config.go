package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const envAPIKey = "BASELINE_API_KEY"

type Config struct {
	APIKey string `json:"api_key"`
}

func Path() (string, error) {
	if path := os.Getenv("BASELINE_CONFIG"); path != "" {
		return path, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "baseline", "config.json"), nil
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode %s: %w", path, err)
	}
	return cfg, nil
}

func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	return os.WriteFile(path, content, 0o600)
}

func APIKey() (string, string, error) {
	if value := os.Getenv(envAPIKey); value != "" {
		return value, envAPIKey, nil
	}

	cfg, err := Load()
	if err != nil {
		return "", "", err
	}
	if cfg.APIKey != "" {
		return cfg.APIKey, "config", nil
	}
	return "", "", nil
}

func SetAPIKey(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("api key cannot be empty")
	}

	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.APIKey = value
	return Save(cfg)
}

func UnsetAPIKey() error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.APIKey = ""
	return Save(cfg)
}

func Mask(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}
