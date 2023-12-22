package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Debug        bool          `json:"debug"`
	ServerConfig *ServerConfig `json:"serverConfig"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func GetConfig(configPath string) (*Config, error) {
	clean := filepath.Clean(configPath)
	configFile, err := os.Open(clean)
	if err != nil {
		return &Config{}, err
	}
	var cfg Config
	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}
