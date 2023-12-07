package service

import (
	"errors"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
)

var (
	ErrNegativeRadius = errors.New("radius must be positive")
)

type MainService struct {
	r repository.DriversRepository
}

func CreateMainService() *MainService {
	return &MainService{r: repository.CreateMapRepository()}
}

func getDistance(lhs repository.Location, rhs repository.Location) float64 {
	return 1
}

func (service MainService) GetDrivers(location repository.Location, radius float64) ([]repository.Driver, error) {
	if radius < 0 {
		return nil, ErrNegativeRadius
	}
	drivers := service.r.GetAllDrivers()
	return drivers, nil
}

func (service MainService) SetDriverLocation(driverId string, location repository.Location) {
	service.r.SetDriverLocation(driverId, location)
}
