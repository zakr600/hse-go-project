package adapter

import (
	"DriverService/internal/adapter/handlers"
	"DriverService/internal/config"
	"github.com/gorilla/mux"
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
) *Server {
	server := &Server{
		log:    log,
		config: cfg,
		router: mux.NewRouter(),
	}
	return server
}

func (s *Server) SetUp() {
	controller := handlers.NewController(s.log)
	trips := s.router.PathPrefix("/trips").Subrouter()

	trips.HandleFunc("", controller.HandlerGetTrips())
	trips.HandleFunc("/{trip_id}", controller.HandlerGetTripByID())
	trips.HandleFunc("/{trip_id}/cancel", controller.HandlerCancelTrip())
	trips.HandleFunc("/{trip_id}/accept", controller.HandlerAcceptTrip())
	trips.HandleFunc("/{trip_id}/start", controller.HandlerStartTrip())

	s.router.HandleFunc("/add", controller.HandlerAddTrip())
}

func (s *Server) Start() error {
	s.log.Info("Starting server")
	return http.ListenAndServe(":"+s.config.ServerConfig.Port, s.router)
}

func (s *Server) Stop() error {
	s.log.Info("Stopping server")
	return nil
}
