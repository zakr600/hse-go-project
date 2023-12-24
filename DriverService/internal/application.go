package internal

import (
	"DriverService/internal/adapter"
	"DriverService/internal/config"
	"DriverService/internal/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type Application struct {
	cfg         *config.Config
	log         *zap.Logger
	server      *adapter.Server
	mongoClient *mongo.Client
}

func NewApplication(cfg *config.Config) *Application {
	log, _ := logger.GetLogger(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	trips := client.Database("driver_service").Collection("trips")

	server := adapter.New(log, cfg, trips)

	if err != nil {
		fmt.Println(err)
	}

	return &Application{cfg: cfg, log: log, server: server, mongoClient: client}
}

func (a *Application) Run(ctx context.Context) error {
	if a.server != nil {
		_ = a.server.Start()
	}
	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if a.mongoClient != nil {
		go func(ctx context.Context) {
			_ = a.mongoClient.Disconnect(ctx)
		}(ctx)
	}
	if a.server != nil {
		_ = a.server.Stop()
	}

	return nil
}
