package inmemory

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/errors"
	"sync"
)

type Repository struct {
	data map[string]models.Trip
	mu   sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]models.Trip),
	}
}

func (repo *Repository) GetAllTrips() ([]models.Trip, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	trips := make([]models.Trip, 0, len(repo.data))
	for _, trip := range repo.data {
		trips = append(trips, trip)
	}
	return trips, nil
}

func (repo *Repository) Get(id string) (*models.Trip, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if value, ok := repo.data[id]; ok {
		return &value, nil
	}
	return nil, errors.NotFoundError{Key: id}
}

func (repo *Repository) Add(value models.Trip) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.data[value.Id]; ok {
		return errors.DuplicateKeyError{Key: value.Id}
	}
	repo.data[value.Id] = value
	return nil
}

func (repo *Repository) ChangeTripStatus(id string, status string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if value, ok := repo.data[id]; ok {
		value.Status = status
		repo.data[id] = value
		return nil
	}
	return errors.NotFoundError{Key: id}
}
