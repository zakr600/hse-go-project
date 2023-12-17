package handlers

import (
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/service"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"net/http"
	"strconv"
)

var (
	DriverIdParam = "driver_id"
	LAT           = "lat"
	LNG           = "lng"
	RADIUS        = "radius"
	TracerName    = "driver_service"
	BitSize       = 64
)

type Controller struct {
	service            *service.MainService
	getDriversCounter  prometheus.Counter
	setLocationCounter prometheus.Counter
}

var tracer = otel.Tracer(TracerName)

func NewController(service *service.MainService) *Controller {
	return &Controller{
		service: service,
		getDriversCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "adapter", Name: "get_drivers_count", Help: "Get drivers controller usage counter",
		}),
		setLocationCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "adapter", Name: "set_location_count", Help: "Set location controller usage counter",
		}),
	}
}

func (controller *Controller) GetDrivers(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "/drivers")
	defer span.End()

	controller.getDriversCounter.Inc()

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	lat, err := strconv.ParseFloat(r.Form.Get(LAT), BitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	lng, err := strconv.ParseFloat(r.Form.Get(LNG), BitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	radius, err := strconv.ParseFloat(r.Form.Get(RADIUS), BitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	drivers, err := controller.service.GetDrivers(repository.CreateLocation(lat, lng), radius)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(drivers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	type driverDTO struct {
		Lat      float64 `json:"lat"`
		Lng      float64 `json:"lng"`
		DriverId string  `json:"id"`
	}
	driversDTO := make([]driverDTO, 0)
	for _, driver := range drivers {
		driversDTO = append(driversDTO, driverDTO{Lat: driver.Location.Lat, Lng: driver.Location.Lng, DriverId: driver.DriverId})
	}
	utils.SendJson(w, r, driversDTO, http.StatusOK)
}

func (controller *Controller) SetDriverLocation(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "/driver/{}/location")
	defer span.End()

	controller.setLocationCounter.Inc()

	driverId := chi.URLParam(r, DriverIdParam)

	var locationDTO struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if !utils.TryParseJsonQuery(w, r, &locationDTO) {
		return
	}

	controller.service.SetDriverLocation(driverId, repository.CreateLocation(locationDTO.Lat, locationDTO.Lng))
	w.WriteHeader(http.StatusOK)
}
