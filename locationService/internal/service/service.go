package service

import (
	"errors"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/utils"
)

var (
	ErrNegativeRadius = errors.New("radius must be positive")
)

type MainService struct {
	r repository.DriversRepository
}

// CreateMainService creates main service for application
func CreateMainService(repo repository.DriversRepository) *MainService {
	return &MainService{r: repo}
}

// GetDrivers gets list of drivers in selected radius
func (service MainService) GetDrivers(location repository.Location, radius float64) ([]repository.Driver, error) {
	if radius < 0 {
		return nil, ErrNegativeRadius
	}
	allDrivers := service.r.GetAllDrivers()

	selectedDrivers := make([]repository.Driver, 0)
	for _, driver := range allDrivers {
		if utils.GetDistance(location, driver.Location) <= radius {
			selectedDrivers = append(selectedDrivers, driver)
		}
	}
	return selectedDrivers, nil
}

// SetDriverLocation sets driver location
func (service MainService) SetDriverLocation(driverId string, location repository.Location) {
	service.r.SetDriverLocation(driverId, location)
}
