package main

import (
	"context"
	"fmt"
	"log"

	"DriverService/internal"
	"DriverService/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err.Error())
	}

	fmt.Println(*cfg)
	fmt.Println(cfg.Debug)
	fmt.Println(cfg.ServerConfig.Port)
	fmt.Println(cfg.ServerConfig.Host)

	app, err := createApplication(cfg)
	if err != nil {
		log.Fatal("Failed to create the application: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Println("Application returned with error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Received signal: %s\n", <-sigChan)

	timeAfterSignal := 0 * time.Second
	log.Printf("Termination in %s\n", timeAfterSignal)
	time.Sleep(timeAfterSignal)

	cancel()

	log.Println("Application terminated")
}

func loadConfig(configPath string) (*config.Config, error) {
	cfg, err := config.GetConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func createApplication(cfg *config.Config) (*internal.Application, error) {
	app := internal.NewApplication(cfg)

	return app, nil
}
