package handlers

import (
	"DriverService/internal/models"
	"DriverService/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(trips); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (controller *Controller) HandlerGetTripByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tripID := vars["trip_id"]
		controller.log.Info("Request: get_trip %s", zap.String("trip_id", tripID))
		trip, err := controller.s.GetTrip(tripID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		controller.log.Info("Request: cancel trip  %s", zap.String("trip_id", tripID))
		err := controller.s.SetTripStatus(tripID, "CANCELED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

func (controller *Controller) HandlerAcceptTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tripID := vars["trip_id"]

		controller.log.Info("Request: accept trip  %s", zap.String("trip_id", tripID))
		err := controller.s.SetTripStatus(tripID, "DRIVER_FOUND")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

func (controller *Controller) HandlerStartTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tripID := vars["trip_id"]

		controller.log.Info("Request: start trip  %s", zap.String("trip_id", tripID))
		err := controller.s.SetTripStatus(tripID, "STARTED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// DEBUG
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
