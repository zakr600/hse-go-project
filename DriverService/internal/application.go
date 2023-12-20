package internal

import (
	"DriverService/internal/adapter"
	"DriverService/internal/config"
	"DriverService/internal/logger"
	"context"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	log    *zap.Logger
	server *adapter.Server
}

func NewApplication(cfg *config.Config) *Application {
	log, _ := logger.GetLogger(cfg)
	server := adapter.New(log, cfg)
	server.SetUp()
	return &Application{cfg: cfg, log: log, server: server}
}

func (a *Application) Run(ctx context.Context) error {
	if a.server != nil {
		_ = a.server.Start()
	}
	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	if a.server != nil {
		_ = a.server.Stop()
	}
	return nil
}
