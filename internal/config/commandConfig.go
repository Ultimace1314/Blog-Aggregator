package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

func readConfig() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homePath + "/" + configFileName, nil
}

func writeConfigFile(cfg Config) error {
	return nil
}
