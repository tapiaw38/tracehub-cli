package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the CLI configuration
type Config struct {
	CurrentProject *ProjectConfig `json:"current_project,omitempty"`
	ServerURL      string          `json:"server_url"`
	APIKey         string          `json:"api_key"`
}

// ProjectConfig represents a configured project
type ProjectConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var configPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	configPath = filepath.Join(home, ".tracehub", "config.json")
}

// Load loads the configuration
func Load() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration
func Save(cfg *Config) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetConfigPath returns the configuration file path
func GetConfigPath() string {
	return configPath
}
