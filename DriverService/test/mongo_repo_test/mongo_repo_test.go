package mongo_repo

import (
	"DriverService/internal"
	"DriverService/internal/config"
	"DriverService/internal/models"
	"DriverService/internal/repository/mongo_db"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"testing"
)

var (
	testPort = "8082"
)

func GenerateTrips() []models.Trip {
	trips := make([]models.Trip, 3)
	for i := range trips {
		trips[i] = models.Trip{
			Id:       fmt.Sprintf("trip%d", i+1),
			DriverId: fmt.Sprintf("driver%d", rand.Intn(10)+1),
			From: &models.LatLngLiteral{
				Lat: rand.Float64()*180 - 90,
				Lng: rand.Float64()*360 - 180,
			},
			To: &models.LatLngLiteral{
				Lat: rand.Float64()*180 - 90,
				Lng: rand.Float64()*360 - 180,
			},
			Price: models.Money{
				Amount:   rand.Float64()*100 + 10,
				Currency: "USD",
			},
			Status: "completed",
		}
	}
	return trips
}

func Setup(t *testing.T) (*mongo.Collection, context.CancelFunc) {
	cfg, _ := config.GetConfig(true)
	cfg.ServerConfig.Port = testPort
	ctx, cancel := context.WithCancel(context.Background())
	app := internal.NewApplication(cfg)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Println("Connecting failed:", err.Error())
	}
	tripsDb := client.Database("driver_service").Collection("trips")

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Println("Application returned with error:", err.Error())
		}
	}()
	return tripsDb, cancel
}

func TestRepo(t *testing.T) {
	sampleTrips := GenerateTrips()
	tripsDb, cancel := Setup(t)
	defer cancel()
	repo, err := mongo_db.NewRepository(tripsDb)
	if err != nil {
		log.Println("Error when creating repo: ", err.Error())
	}

	t.Run("Test insert", func(t *testing.T) {
		t.Run("Test simple", func(t *testing.T) {
			err := repo.Insert(sampleTrips[0])
			if err != nil {
				t.Errorf(err.Error())
			}
			err = repo.Insert(sampleTrips[1])
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	})
}
