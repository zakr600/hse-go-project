package main

import (
	"context"
	"flag"
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
	var cfgPath string
	flag.StringVar(&cfgPath, "cfg", "./configs/config.yaml", "server config")
	flag.Parse()
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		fmt.Println("Failed to read config")
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	serverApp := internal.NewApplication(*cfg, service.CreateMainService(repository.CreateMapRepository()))
	internal.Run(serverApp)
	<-ctx.Done()
	ctx, stop = context.WithTimeout(ctx, 10*time.Second)
	defer stop()
	internal.Stop(serverApp, ctx)
}
