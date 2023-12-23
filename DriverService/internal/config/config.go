package config

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
)

var (
	DriverServerPort        = "DRIVER_SERVER_PORT"
	DriverServerHost        = "DRIVER_SERVER_HOST"
	Debug                   = "DRIVER_SERVICE_DEBUG"
	DefaultDriverServerPort = "8081"
	DefaultDriverServerHost = "localhost"
)

type Config struct {
	Debug        bool          `json:"debug"`
	ServerConfig *ServerConfig `json:"serverConfig"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func GetConfigFromFile(configPath string) (*Config, error) {
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

func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return &Config{}, errors.New("failed to load env")
	}
	serverPort := GetEnvString(DriverServerPort, DefaultDriverServerPort)
	serverHost := GetEnvString(DriverServerHost, DefaultDriverServerHost)

	cfg := &ServerConfig{
		Host: serverHost,
		Port: serverPort,
	}

	debug := GetEnvBool(Debug, false)
	return &Config{
		Debug:        debug,
		ServerConfig: cfg,
	}, nil
}

func GetEnvBool(key string, defaultVal bool) bool {
	envVar, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	val, err := strconv.ParseBool(envVar)
	if err != nil {
		return defaultVal
	}
	return val
}

func GetEnvString(key string, defaultVal string) string {
	envVar, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return envVar
}
