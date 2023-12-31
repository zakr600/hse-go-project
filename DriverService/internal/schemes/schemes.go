package schemes

import (
	"DriverService/internal/models"
	"github.com/google/uuid"
	"time"
)

const (
	AcceptType   = "trip.command.accept"
	StartType    = "trip.command.start"
	EndType      = "trip.command.end"
	CancelType   = "trip.command.cancel"
	driverSource = "/driver"
	appJson      = "application/json"
)

type Command struct {
	ID              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	DataContentType string                 `json:"datacontenttype"`
	Time            time.Time              `json:"time"`
	Data            map[string]interface{} `json:"data"`
}

func NewCommand(t string, data map[string]interface{}) Command {
	return Command{
		ID:              uuid.New().String(),
		Source:          driverSource,
		Type:            t,
		DataContentType: appJson,
		Time:            time.Now(),
		Data:            data,
	}
}

type JsonData struct {
	Event *Event `json:"data"`
}

type Event struct {
	TripID  string                `json:"trip_id"`
	OfferID string                `json:"offer_id"`
	Price   *models.Money         `json:"price"`
	From    *models.LatLngLiteral `json:"from"`
	To      *models.LatLngLiteral `json:"to"`
	Status  string                `json:"status"`
}

func EventToTrip(e *Event) models.Trip {
	return models.Trip{
		Id:       e.TripID,
		DriverId: "",
		From:     e.From,
		To:       e.To,
		Price:    e.Price,
		Status:   e.Status,
	}
}
