package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config represents the persisted application configuration.
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read loads config from the user's home directory.
func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// SetUser sets current_user_name and persists the updated config.
func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, configFileName), nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}
