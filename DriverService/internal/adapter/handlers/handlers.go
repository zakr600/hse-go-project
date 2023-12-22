package handlers

import (
	"DriverService/internal/models"
	"DriverService/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func NewController(log *zap.Logger) *Controller {
	return &Controller{
		s:   service.New(),
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
		controller.log.Info("Request: get_trip %s", zap.String("trip_id", tripID))
		trip, err := controller.s.GetTrip(tripID)
		httpRequestsTotal.WithLabelValues("HandlerGetTripByID").Inc()
		if err != nil {
			httpRequests5xx.WithLabelValues("HandlerGetTripByID").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		httpRequests2xx.WithLabelValues("HandlerGetTripByID").Inc()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(trip); err != nil {
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
		controller.log.Info("Request: cancel trip  %s", zap.String("trip_id", tripID))
		err := controller.s.SetTripStatus(tripID, "CANCELED")
		if err != nil {
			httpRequests5xx.WithLabelValues("HandlerCancelTrip").Inc()
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
		controller.log.Info("Request: accept trip  %s", zap.String("trip_id", tripID))
		err := controller.s.OnStatusAccept(tripID)
		if err != nil {
			httpRequests5xx.WithLabelValues("HandlerAcceptTrip").Inc()
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
		controller.log.Info("Request: start trip  %s", zap.String("trip_id", tripID))
		err := controller.s.SetTripStatus(tripID, "STARTED")
		if err != nil {
			httpRequests5xx.WithLabelValues("HandlerStartTrip").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpRequests2xx.WithLabelValues("HandlerStartTrip").Inc()
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
