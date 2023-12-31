package service

import (
	"DriverService/internal/config"
	"DriverService/internal/models"
	"DriverService/internal/repository"
	"DriverService/internal/repository/mongo_db"
	"DriverService/internal/schemes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var (
	radius = 10.0
)

type Service struct {
	repo   repository.Repository
	writer *kafka.Writer
	reader *kafka.Reader
	cfg    *config.Config
	log    *zap.Logger
}

func New(cfg *config.Config, tripsDb *mongo.Collection, log *zap.Logger) *Service {
	kafkaAddress := os.Getenv("KAFKA")
	repo, err := mongo_db.NewRepository(tripsDb)

	if err != nil {
		log.Error(err.Error())
	}
	return &Service{
		repo: repo,
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
		log: log,
	}
}

func (s *Service) GetTrips() ([]models.Trip, error) {
	trips, err := s.repo.GetAllTrips()
	return trips, err
}

func (s *Service) GetTrip(tripId string) (*models.Trip, error) {
	trip, err := s.repo.Get(tripId)
	return trip, err
}

func (s *Service) OnAcceptTrip(tripID string) error {
	err := s.repo.SetStatus(tripID, "ACCEPTED")
	if err != nil {
		return err
	}
	trip, err := s.repo.Get(tripID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"trip_id":   tripID,
		"driver_id": trip.DriverId,
	}
	command := schemes.NewCommand(schemes.AcceptType, data)

	err = s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnStartTrip(tripID string) error {
	err := s.repo.SetStatus(tripID, "STARTED")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"trip_id": tripID,
	}
	command := schemes.NewCommand(schemes.AcceptType, data)

	err = s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnEndTrip(tripID string) error {
	err := s.repo.SetStatus(tripID, "ENDED")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"trip_id": tripID,
	}
	command := schemes.NewCommand(schemes.AcceptType, data)

	err = s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnCancelTrip(tripID string) error {
	err := s.repo.SetStatus(tripID, "CANCELLED")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"trip_id": tripID,
		"reason":  "Cancelled",
	}
	command := schemes.NewCommand(schemes.AcceptType, data)

	err = s.writeCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OnCreateTrip(event *schemes.Event) error {
	s.log.Debug("Trip Created")
	trip := schemes.EventToTrip(event)
	drivers, err := s.GetDrivers(trip.To.Lat, trip.To.Lng, radius)
	if err != nil {
		return err
	}
	for _, driver := range drivers {
		curTrip := trip
		curTrip.DriverId = driver.DriverId
		err = s.repo.Insert(curTrip)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) AddTrip(trip models.Trip) error {
	err := s.repo.Insert(trip)
	if err != nil {
		s.log.Error(err.Error())
	}
	return nil
}

func (s *Service) writeCommand(cmd schemes.Command) error {
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
			s.log.Error(err.Error())
			break
		}
		var jsonData schemes.JsonData
		err = json.Unmarshal(m.Value, &jsonData)
		if err != nil {
			s.log.Error(err.Error())
			continue
		}
		err = s.OnCreateTrip(jsonData.Event)
		if err != nil {
			s.log.Error(err.Error())
			continue
		}
	}

	if err := s.reader.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetDrivers(lat, lng, radius float64) ([]models.Driver, error) {
	url := fmt.Sprintf("http://%s:%s/drivers?lat=%f&lng=%f&radius=%f", s.cfg.LocationServiceConfig.Host, s.cfg.LocationServiceConfig.Port, lat, lng, radius)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("drivers not found")
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unexpected response status: %d", resp.StatusCode))
	}

	var drivers []models.Driver
	if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}
