package service

import (
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/service"
	"testing"
)

func TestService(t *testing.T) {
	serv := service.CreateMainService(repository.CreateMapRepository())

	serv.SetDriverLocation("1", repository.CreateLocation(0, 0))

	drivers, err := serv.GetDrivers(repository.CreateLocation(0, 0), 1)
	if len(drivers) != 1 || err != nil || drivers[0].DriverId != "1" || drivers[0].Location.Lat != 0 || drivers[0].Location.Lng != 0 {
		t.Errorf("Not found driver or err != nil")
	}

	drivers, err = serv.GetDrivers(repository.CreateLocation(100, 100), 0.1)
	if len(drivers) != 0 || err != nil {
		t.Errorf("Not found driver or err != nil")
	}

	drivers, err = serv.GetDrivers(repository.CreateLocation(0, 0), -10)
	if err == nil {
		t.Errorf("Negative radius exception not thrown")
	}
}
