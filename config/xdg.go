package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func basePath() string {
	xdgConfigDir, ok := os.LookupEnv("XDG_CONFIG_HOME")

	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			home = os.Getenv("HOME")
		}
		xdgConfigDir = filepath.Join(home, ".config")
	}

	return filepath.Join(xdgConfigDir, "swoosh")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func xdgConfigFile(defaultConfig Config) (string, error) {
	configDir := basePath()

	if !exists(configDir) {
		err := os.Mkdir(configDir, 0o755)
		if err != nil {
			return "", fmt.Errorf("failed to create config dir: %v", err)
		}
	}

	configFile := filepath.Join(configDir, "config.json")

	if !exists(configFile) {
		jsonData, err := json.MarshalIndent(defaultConfig, "", " ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal default config file: %v", err)
		}

		err = os.WriteFile(configFile, jsonData, 0o644)
		if err != nil {
			return "", fmt.Errorf("failed to write config file: %v", err)
		}
	}

	return configFile, nil
}
