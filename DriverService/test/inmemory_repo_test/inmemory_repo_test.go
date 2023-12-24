package inmemory_repo

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/inmemory"
	"DriverService/internal/trip_errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	//"DriverService/internal/repository/inmemory"
)

func GenerateTrips() []models.Trip {
	trips := make([]models.Trip, 3)
	for i := range trips {
		trips[i] = models.Trip{
			Id:       fmt.Sprintf("trip%d", i+1),
			DriverId: fmt.Sprintf("driver%d", rand.Intn(10)+1),
			From: &models.LatLngLiteral{
				Lat: rand.Float64()*180 - 90,
				Lng: rand.Float64()*360 - 180,
			},
			To: &models.LatLngLiteral{
				Lat: rand.Float64()*180 - 90,
				Lng: rand.Float64()*360 - 180,
			},
			Price: models.Money{
				Amount:   rand.Float64()*100 + 10,
				Currency: "USD",
			},
			Status: "completed",
		}
	}
	return trips
}

func TestRepo(t *testing.T) {
	repo := inmemory.NewRepository()
	sampleTrips := GenerateTrips()

	t.Run("Test insert", func(t *testing.T) {
		t.Run("Test simple", func(t *testing.T) {
			err := repo.Insert(sampleTrips[0])
			if err != nil {
				t.Errorf(err.Error())
			}
			err = repo.Insert(sampleTrips[1])
			if err != nil {
				t.Errorf(err.Error())
			}
		})
		t.Run("Test duplicate key", func(t *testing.T) {
			err := repo.Insert(sampleTrips[0])
			err = repo.Insert(sampleTrips[0])
			if reflect.TypeOf(err) != reflect.TypeOf(trip_errors.DuplicateKeyError{}) {
				t.Errorf("Expected DuplicateKeyError, got %s", err)
			}
		})
	})
	t.Run("Test get", func(t *testing.T) {
		_ = repo.Insert(sampleTrips[0])
		_ = repo.Insert(sampleTrips[1])
		t.Run("Test Not Found", func(t *testing.T) {
			_, err := repo.Get("abc")
			if reflect.TypeOf(err) != reflect.TypeOf(trip_errors.NotFoundError{}) {
				t.Errorf("Expected NotFoundError, got %s", err)
			}
		})
		t.Run("Test Get", func(t *testing.T) {
			trip, err := repo.Get(sampleTrips[0].Id)
			if err != nil {
				t.Errorf(err.Error())
			}
			if *trip != sampleTrips[0] {
				t.Errorf("Got incorrect trip %v", trip)
			}
		})
	})
	t.Run("Test set status", func(t *testing.T) {
		_ = repo.Insert(sampleTrips[0])
		_ = repo.Insert(sampleTrips[1])
		t.Run("Test Not Found", func(t *testing.T) {
			err := repo.SetStatus("abc", "STARTED")
			if reflect.TypeOf(err) != reflect.TypeOf(trip_errors.NotFoundError{}) {
				t.Errorf("Expected NotFoundError, got %s", err)
			}
		})
		t.Run("Test change trip", func(t *testing.T) {
			err := repo.SetStatus(sampleTrips[0].Id, "STARTED")
			if err != nil {
				t.Errorf(err.Error())
			}
			trip, err := repo.Get(sampleTrips[0].Id)
			if trip.Status != "STARTED" {
				t.Errorf("Expected status STARTED, got %s", trip.Status)
			}
		})
	})
}
