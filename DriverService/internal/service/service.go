package service

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/inmemory"
	"DriverService/internal/schemes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
)

type Service struct {
	repo   *inmemory.Repository
	writer *kafka.Writer
	reader *kafka.Reader
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
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{kafkaAddress},
			GroupID:  "consumer-group-id",
			Topic:    "events",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
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

func (s *Service) OnAcceptTrip(tripID string) error {
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

func (s *Service) OnStartTrip(tripID string) error {
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

func (s *Service) OnEndTrip(tripID string) error {
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

func (s *Service) OnCancelTrip(tripID string) error {
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

func (s *Service) OnCreateTrip(event schemes.Event) error {
	fmt.Println("Trip created")
	fmt.Println("Event", event)
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

func (s *Service) FetchEvents() error {
	for {
		m, err := s.reader.ReadMessage(context.Background())
		if err != nil {
			// log ERROR
			break
		}
		var jsonData schemes.JsonData
		err = json.Unmarshal(m.Value, &jsonData)
		if err != nil {
			// log error
			continue
		}
		err = s.OnCreateTrip(jsonData.Event)
		if err != nil {
			// log error
			continue
		}
	}

	if err := s.reader.Close(); err != nil {
		return err
	}
	return nil
}
