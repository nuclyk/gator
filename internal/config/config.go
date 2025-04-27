package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_URL          string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func write(cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0777)
	return nil

}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path += "/" + configFileName

	return path, nil
}

func Read(fileName string) (Config, error) {
	path, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	var config Config
	json.Unmarshal(data, &config)
	return config, nil
}
