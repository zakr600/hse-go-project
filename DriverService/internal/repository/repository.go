package repository

import (
	"DriverService/internal/models"
	"DriverService/internal/repository/inmemory"
	"DriverService/internal/repository/mongo_db"
)

var _ Repository = &mongo_db.Repository{}
var _ Repository = &inmemory.Repository{}

type Repository interface {
	GetAllTrips() ([]models.Trip, error)
	Get(id string) (*models.Trip, error)
	Insert(value models.Trip) error
	SetStatus(id string, status string) error
}
