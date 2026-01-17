package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GamesPath     string `json:"games_path"`
	ShareName     string `json:"share_name"`
	UseGuest      bool   `json:"use_guest"`
	SambaUser     string `json:"samba_user,omitempty"`
	ConfigVersion string `json:"config_version"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "ps2smb")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found, run 'ps2smb init' first")
		}
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

// Exists checks if config file exists
func Exists() bool {
	configPath, err := GetConfigPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(configPath)
	return err == nil
}
