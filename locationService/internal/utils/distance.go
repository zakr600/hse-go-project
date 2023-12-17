package utils

import (
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"github.com/umahmood/haversine"
)

func toHaversine(location repository.Location) haversine.Coord {
	return haversine.Coord{Lat: location.Lat, Lon: location.Lng}
}

func GetDistance(lhs repository.Location, rhs repository.Location) float64 {
	_, km := haversine.Distance(toHaversine(lhs), toHaversine(rhs))
	return km
}
