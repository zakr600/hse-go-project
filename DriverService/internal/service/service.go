package service

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/inmemory"
	"context"
	"fmt"
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

func (s *Service) SetTripStatus(tripId string, status string) error {
	return s.repo.ChangeTripStatus(tripId, status)
}

func (s *Service) AddTrip(trip models.Trip) {
	s.repo.Add(trip)
}

func (s *Service) OnStatusAccept(tripID string) error {
	ctx := context.Background()
	_ = s.repo.ChangeTripStatus(tripID, "ACCEPTED")
	fmt.Println(os.Getenv("KAFKA"))
	err := s.writer.WriteMessages(ctx, kafka.Message{Key: []byte("command"), Value: []byte("abc")})
	if err != nil {
		fmt.Println("KAFKA ERROR", err.Error())
	}
	return nil
}
