package repository

import (
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"testing"
)

func TestMapRepository(t *testing.T) {
	repo := repository.CreateMapRepository()

	repo.SetDriverLocation("0", repository.CreateLocation(0.2, 1.3))

	drivers := repo.GetAllDrivers()
	if len(drivers) != 1 {
		t.Errorf("Drivers count inconsisntent")
	}
	driver := drivers[0]
	if driver.Location.Lat != 0.2 || driver.Location.Lng != 1.3 {
		t.Errorf("Driver location inconsistent")
	}
}
