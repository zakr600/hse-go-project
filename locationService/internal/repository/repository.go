package repository

type Location struct {
	Lat float64
	Lng float64
}

type Driver struct {
	DriverId string
	Location Location
}

func CreateLocation(lat float64, lng float64) Location {
	return Location{Lat: lat, Lng: lng}
}

func CreateDriver(driverId string, location Location) Driver {
	return Driver{DriverId: driverId, Location: location}
}

type DriversRepository interface {
	SetDriverLocation(string, Location)
	GetAllDrivers() []Driver
}
