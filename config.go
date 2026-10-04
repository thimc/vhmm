package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	GameDir       string `json:"game_dir"`
	RepositoryURL string `json:"repository_url"`
	Script        string `json:"script"`
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile("config.json")
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
