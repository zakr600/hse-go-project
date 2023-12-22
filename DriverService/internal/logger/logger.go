package logger

import (
	"DriverService/internal/config"
	"go.uber.org/zap"
)

func GetLogger(config *config.Config) (*zap.Logger, error) {
	if config.Debug {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}
