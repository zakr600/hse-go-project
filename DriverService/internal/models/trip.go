package models

type Trip struct {
	Id       string         `json:"id"`
	DriverId string         `json:"driverId"`
	From     *LatLngLiteral `json:"from"`
	To       *LatLngLiteral `json:"to"`
	Price    Money          `json:"price"`
	Status   string         `json:"status"`
}
