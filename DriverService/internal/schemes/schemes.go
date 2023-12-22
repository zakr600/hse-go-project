package schemes

import (
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

type Scheme struct {
	ID              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	DataContentType string                 `json:"datacontenttype"`
	Time            time.Time              `json:"time"`
	Data            map[string]interface{} `json:"data"`
}

func NewScheme(t string, data map[string]interface{}) Scheme {
	return Scheme{
		ID:              uuid.New().String(),
		Source:          driverSource,
		Type:            t,
		DataContentType: appJson,
		Time:            time.Now(),
		Data:            data,
	}
}
