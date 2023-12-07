package handlers

import (
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var (
	DriverIdParam = "driver_id"
)

type Controller struct {
	service *service.MainService
}

func NewController(service *service.MainService) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller *Controller) GetDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := controller.service.GetDrivers(repository.CreateLocation(0, 0), 1)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for _, driver := range drivers {
		w.Write([]byte(driver.DriverId))
	}

	w.WriteHeader(200)
}

func (controller *Controller) SetDriverLocation(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(chi.URLParam(r, DriverIdParam)))
	if err != nil {
		return
	}
	driverId := chi.URLParam(r, DriverIdParam)
	controller.service.SetDriverLocation(driverId, repository.CreateLocation(0, 0))
	w.WriteHeader(200)
}
