package main

import (
	"encoding/json"
	"os"
)

const configFile = "config.json"

type Config struct {
	GameDir       string `json:"game_dir"`
	RepositoryURL string `json:"repository_url"`
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
