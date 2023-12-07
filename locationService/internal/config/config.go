package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type ServerConfig struct {
	Port       int64  `yaml:"port"`
	ApiVersion string `yaml:"version"`
	Debug      bool   `yaml:"debug"`
}

func NewConfig(filePath string) (*ServerConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
