package main

import (
	"context"
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"os/signal"
	"syscall"
	"time"

	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/config"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Failed to read config")
		return
	}
	fmt.Println(cfg.Port)
	fmt.Println(cfg.ApiVersion)
	fmt.Println(cfg.Debug)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	serverApp := internal.NewApplication(*cfg, service.CreateMainService(repository.CreateMapRepository()))
	internal.Run(serverApp)
	<-ctx.Done()
	ctx, stop = context.WithTimeout(ctx, 10*time.Second)
	defer stop()
	internal.Stop(serverApp, ctx)
}
