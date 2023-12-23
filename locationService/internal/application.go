package internal

import (
	"context"
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/config"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/handlers"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/middleware"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	cfg    *config.ServerConfig
	server *http.Server
	logger *zap.Logger
}

func (app *App) Run() {
	app.logger.Info("Starting app....")
	go func() {
		err := app.server.ListenAndServe()
		if err != nil {
			app.logger.Error(err.Error())
		}
	}()
}

func (app *App) Stop(ctx context.Context) {
	app.logger.Info("Stopping app....")
	err := app.server.Shutdown(ctx)
	if err != nil {
		app.logger.Error("Failed to stop app!")
		app.logger.Error(err.Error())
	} else {
		app.logger.Info("App stopped successfully")
	}
}

func getController(serv *service.MainService) http.Handler {
	r := chi.NewRouter()
	controller := handlers.NewController(serv)
	r.Use(middleware.Logger)
	r.Get("/drivers", controller.GetDrivers)
	r.Post("/drivers/{driver_id}/location", controller.SetDriverLocation)
	r.Handle("/metrics", promhttp.Handler())
	return r
}

func NewApplication(config config.ServerConfig, service *service.MainService) *App {
	log, err := logging.GetLogger(config.Debug)
	if err != nil {
		fmt.Println("Failed to create logger!")
		panic(err)
	}
	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: getController(service),
		},
		cfg:    &config,
		logger: log,
	}
}

func Run(app *App) {
	app.Run()
}

func Stop(app *App, ctx context.Context) {
	app.Stop(ctx)
}
