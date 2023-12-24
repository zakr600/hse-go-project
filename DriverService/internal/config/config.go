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
	DriverServerPort           = "DRIVER_SERVER_PORT"
	DriverServerHost           = "DRIVER_SERVER_HOST"
	MongoURI                   = "MONGO_URI"
	MongoMigrationsPath        = "MONGO_MIGRATIONS_PATH"
	Debug                      = "DRIVER_SERVICE_DEBUG"
	LocationServicePort        = "LOCATION_SERVICE_PORT"
	LocationServiceHost        = "LOCATION_SERVICE_HOST"
	DefaultDriverServerPort    = "8081"
	DefaultDriverServerHost    = "localhost"
	DefaultMongoURI            = "mongodb://mongodb:27017"
	DefaultMongoMigrationsPath = "./migrations"
	DefaultLocationServicePort = "8080"
	DefaultLocationServiceHost = "localhost"
)

type Config struct {
	Debug                 bool                   `json:"debug"`
	ServerConfig          *ServerConfig          `json:"serverConfig"`
	MongoURI              string                 `json:"mongoUri"`
	MongoMigrationsPath   string                 `json:"mongoMigrationsPath"`
	LocationServiceConfig *LocationServiceConfig `json:locationServiceConfig`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type LocationServiceConfig struct {
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

func GetConfig(useDefaults bool) (*Config, error) {
	err := godotenv.Load()
	if !useDefaults && err != nil {
		return &Config{}, errors.New("failed to load env")
	}
	serverPort := GetEnvString(DriverServerPort, DefaultDriverServerPort)
	serverHost := GetEnvString(DriverServerHost, DefaultDriverServerHost)
	MongoURI := GetEnvString(MongoURI, DefaultMongoURI)
	MongoMigrationsPath := GetEnvString(MongoMigrationsPath, DefaultMongoMigrationsPath)

	serverCfg := &ServerConfig{
		Host: serverHost,
		Port: serverPort,
	}

	LocationServiceHost := GetEnvString(LocationServiceHost, DefaultLocationServiceHost)
	LocationServicePort := GetEnvString(LocationServicePort, DefaultLocationServicePort)

	locationCfg := &LocationServiceConfig{
		Host: LocationServiceHost,
		Port: LocationServicePort,
	}

	debug := GetEnvBool(Debug, false)
	return &Config{
		Debug:                 debug,
		ServerConfig:          serverCfg,
		MongoURI:              MongoURI,
		MongoMigrationsPath:   MongoMigrationsPath,
		LocationServiceConfig: locationCfg,
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
