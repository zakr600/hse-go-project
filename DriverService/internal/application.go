package internal

import (
	"DriverService/internal/adapter"
	"DriverService/internal/config"
	"DriverService/internal/logger"
	"DriverService/internal/migrations"
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

	migration := migrations.NewMigration(client, client.Database("driver_service"), log)
	err = migration.Run(cfg.MongoMigrationsPath)
	if err != nil {
		fmt.Println(err)
	}

	server := adapter.New(log, cfg, trips)

	if err != nil {
		fmt.Println(err)
	}

	return &Application{cfg: cfg, log: log, server: server, mongoClient: client}
}

func (a *Application) Run(ctx context.Context) error {
	a.log.Info("Starting Application...")
	if a.server != nil {
		err := a.server.Start()
		if err != nil {
			a.log.Error(fmt.Sprintf("Couldn't start Server: %s", err.Error()))
		}
	}
	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	a.log.Info("Stopping Application...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if a.mongoClient != nil {
		go func(ctx context.Context) {
			err := a.mongoClient.Disconnect(ctx)
			if err != nil {
				a.log.Error(fmt.Sprintf("Mongo client couldn't disconnect: %s", err.Error()))
			}
		}(ctx)
	}
	if a.server != nil {

		go func(ctx context.Context) {
			err := a.server.Stop()
			if err != nil {
				a.log.Error(fmt.Sprintf("Server couldn't stop: %s", err.Error()))
			}
		}(ctx)
	}

	return nil
}
