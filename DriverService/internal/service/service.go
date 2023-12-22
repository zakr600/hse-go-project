package service

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/inmemory"
	"DriverService/internal/schemes"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"os"
)

type Service struct {
	repo   *inmemory.Repository
	writer *kafka.Writer
}

func New() *Service {
	kafkaAddress := os.Getenv("KAFKA")

	return &Service{
		repo: inmemory.NewRepository(),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(kafkaAddress),
			Topic:    "commands",
			Balancer: &kafka.LeastBytes{},
		},
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

func (s *Service) OnStatusAccept(tripID string) error {
	_ = s.repo.ChangeTripStatus(tripID, "ACCEPTED")
	trip, err := s.repo.Get(tripID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"trip_id":   tripID,
		"driver_id": trip.DriverId,
	}
	command := schemes.NewScheme(schemes.AcceptType, data)

	err = s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnStatusStart(tripID string) error {
	_ = s.repo.ChangeTripStatus(tripID, "STARTED")

	data := map[string]interface{}{
		"trip_id": tripID,
	}
	command := schemes.NewScheme(schemes.AcceptType, data)

	err := s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnStatusEnd(tripID string) error {
	_ = s.repo.ChangeTripStatus(tripID, "ENDED")

	data := map[string]interface{}{
		"trip_id": tripID,
	}
	command := schemes.NewScheme(schemes.AcceptType, data)

	err := s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnStatusCancel(tripID string) error {
	_ = s.repo.ChangeTripStatus(tripID, "CANCELLED")

	data := map[string]interface{}{
		"trip_id": tripID,
		"reason":  "Cancelled",
	}
	command := schemes.NewScheme(schemes.AcceptType, data)

	err := s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

// TODO: удалить
func (s *Service) AddTrip(trip models.Trip) {
	s.repo.Add(trip)
}

func (s *Service) writeCommand(cmd schemes.Scheme) error {
	msgBytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: msgBytes,
	}

	return s.writer.WriteMessages(context.Background(), msg)
}
