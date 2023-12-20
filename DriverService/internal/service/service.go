package service

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/inmemory"
)

type Service struct {
	repo *inmemory.Repository
}

func New() *Service {
	return &Service{
		repo: inmemory.NewRepository(),
	}
}

func (s *Service) GetTrips() ([]models.Trip, error) {
	trips := s.repo.GetAllTrips()
	return trips, nil
}

func (s *Service) GetTrip(tripId string) (models.Trip, error) {
	trip, err := s.repo.Get(tripId)
	return trip, err
}

func (s *Service) SetTripStatus(tripId string, status string) error {
	return s.repo.ChangeTripStatus(tripId, status)
}

func (s *Service) AddTrip(trip models.Trip) {
	s.repo.Add(trip)
}
