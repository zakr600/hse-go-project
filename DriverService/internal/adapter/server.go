package adapter

import (
	"DriverService/internal/adapter/handlers"
	"DriverService/internal/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
	log        *zap.Logger
	router     *mux.Router
}

func New(
	log *zap.Logger,
	cfg *config.Config,
	tripsDb *mongo.Collection,
) *Server {
	router := mux.NewRouter()
	SetUp(router, log, tripsDb, cfg)
	server := &Server{
		log:    log,
		config: cfg,
		router: router,
	}
	return server
}

func SetUp(router *mux.Router, log *zap.Logger, tripsDb *mongo.Collection, cfg *config.Config) {
	controller := handlers.NewController(cfg, tripsDb, log)
	log.Info("Registered metrics handler")
	router.Handle("/metrics", promhttp.Handler())
	trips := router.PathPrefix("/trips").Subrouter()

	trips.HandleFunc("", controller.HandlerGetTrips())
	trips.HandleFunc("/{trip_id}", controller.HandlerGetTripByID())
	trips.HandleFunc("/{trip_id}/cancel", controller.HandlerCancelTrip())
	trips.HandleFunc("/{trip_id}/accept", controller.HandlerAcceptTrip())
	trips.HandleFunc("/{trip_id}/start", controller.HandlerStartTrip())
	trips.HandleFunc("/{trip_id}/end", controller.HandlerEndTrip())

	router.HandleFunc("/add", controller.HandlerAddTrip())
}

func (s *Server) Start() error {
	s.log.Info("Starting server")
	return http.ListenAndServe(":"+s.config.ServerConfig.Port, s.router)
}

func (s *Server) Stop() error {
	s.log.Info("Stopping server")
	return nil
}
