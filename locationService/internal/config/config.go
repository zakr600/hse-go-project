package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

var (
	Port           = "port"
	ApiVersion     = "api_version"
	Debug          = "debug"
	errNoParameter = errors.New("failed to read server config")
)

type ServerConfig struct {
	Port       int    `yaml:"port"`
	ApiVersion string `yaml:"version"`
	Debug      bool   `yaml:"debug"`
}

// NewConfig Load server config from environment variables
func NewConfig() (*ServerConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errNoParameter
	}
	return &ServerConfig{
		Port:       getEnvAsInt(Port, 80),
		ApiVersion: getEnv(ApiVersion, ""),
		Debug:      getEnvAsBool(Debug, false),
	}, nil
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
