package handlers

import (
	"DriverService/internal/models"
	"DriverService/internal/service"
	"DriverService/internal/trip_errors"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"handler"})
	httpRequests2xx = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_2xx",
		Help: "Number of request with 2xx status code",
	}, []string{"handler"})
	httpRequests4xx = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_4xx",
		Help: "Number of request with 4xx status code",
	}, []string{"handler"})
	httpRequests5xx = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_5xx",
		Help: "Number of request with 5xx status code",
	}, []string{"handler"})
)

type Controller struct {
	s   *service.Service
	log *zap.Logger
}

func NewController(tripsDb *mongo.Collection, log *zap.Logger) *Controller {
	svc := service.New(tripsDb)

	go func() {
		err := svc.FetchEvents()
		if err != nil {
			log.Error(fmt.Sprint("Error while fetching events", err.Error()))
		}
	}()

	return &Controller{
		s:   svc,
		log: log,
	}
}

func (controller *Controller) HandlerGetTrips() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trips, err := controller.s.GetTrips()
		httpRequestsTotal.WithLabelValues("HandlerGetTrips").Inc()
		if err != nil {
			httpRequests5xx.WithLabelValues("HandlerGetTrips").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(trips); err != nil {
			httpRequests5xx.WithLabelValues("HandlerGetTrips").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpRequests2xx.WithLabelValues("HandlerGetTrips").Inc()
	}
}

func (controller *Controller) HandlerGetTripByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tripID := vars["trip_id"]
		controller.log.Debug("Request: get_trip %s", zap.String("trip_id", tripID))
		if _, err := uuid.Parse(tripID); err != nil {
			httpRequests4xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, fmt.Sprintf("Incorrect trip id: %s", err.Error()), http.StatusBadRequest)
			return
		}

		trip, err := controller.s.GetTrip(tripID)
		httpRequestsTotal.WithLabelValues("HandlerGetTripByID").Inc()

		if errors.Is(err, trip_errors.NotFoundError{}) {
			httpRequests4xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		httpRequests2xx.WithLabelValues("HandlerGetTripByID").Inc()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(*trip); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (controller *Controller) HandlerCancelTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		tripID := vars["trip_id"]
		httpRequestsTotal.WithLabelValues("HandlerCancelTrip").Inc()
		controller.log.Debug("Request: cancel trip  %s", zap.String("trip_id", tripID))
		err := controller.s.OnCancelTrip(tripID)

		if errors.Is(err, trip_errors.NotFoundError{}) {
			httpRequests4xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpRequests2xx.WithLabelValues("HandlerCancelTrip").Inc()
		return
	}
}

func (controller *Controller) HandlerAcceptTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tripID := vars["trip_id"]
		httpRequestsTotal.WithLabelValues("HandlerAcceptTrip").Inc()
		controller.log.Debug("Request: accept trip  %s", zap.String("trip_id", tripID))
		err := controller.s.OnAcceptTrip(tripID)
		if errors.Is(err, trip_errors.NotFoundError{}) {
			httpRequests4xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpRequests2xx.WithLabelValues("HandlerAcceptTrip").Inc()
		return
	}
}

func (controller *Controller) HandlerStartTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tripID := vars["trip_id"]

		httpRequestsTotal.WithLabelValues("HandlerStartTrip").Inc()
		controller.log.Debug("Request: start trip  %s", zap.String("trip_id", tripID))
		err := controller.s.OnStartTrip(tripID)
		if errors.Is(err, trip_errors.NotFoundError{}) {
			httpRequests4xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpRequests2xx.WithLabelValues("HandlerStartTrip").Inc()
		return
	}
}

func (controller *Controller) HandlerEndTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tripID := vars["trip_id"]

		httpRequestsTotal.WithLabelValues("HandlerEndTrip").Inc()
		controller.log.Debug("Request: end trip  %s", zap.String("trip_id", tripID))
		err := controller.s.OnEndTrip(tripID)
		if errors.Is(err, trip_errors.NotFoundError{}) {
			httpRequests4xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpRequests2xx.WithLabelValues("HandlerEndTrip").Inc()
		return
	}
}

// HandlerAddTrip DEBUG
func (controller *Controller) HandlerAddTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trip := models.Trip{
			Id:       fmt.Sprintf("%v", rand.Int()),
			DriverId: fmt.Sprintf("%v", rand.Int()),
		}
		controller.s.AddTrip(trip)
		_, _ = io.WriteString(w, "Trip added")
	}
}
