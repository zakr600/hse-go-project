package models

type Driver struct {
	DriverId      string        `json:"driverId"`
	LatLngLiteral LatLngLiteral `json:"location"`
}
