package repository

import (
	"sync"
)

type MapRepository struct {
	mx        sync.Mutex
	locations map[string]Location
}

func CreateMapRepository() *MapRepository {
	return &MapRepository{
		locations: make(map[string]Location),
	}
}

func (r *MapRepository) SetDriverLocation(driverId string, location Location) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.locations[driverId] = location
}

func (r *MapRepository) GetAllDrivers() []Driver {
	r.mx.Lock()
	defer r.mx.Unlock()

	drivers := make([]Driver, 0)
	for driverId, location := range r.locations {
		drivers = append(drivers, CreateDriver(driverId, location))
	}
	return drivers
}
