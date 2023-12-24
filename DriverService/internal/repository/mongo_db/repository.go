package mongo_db

import (
	"DriverService/internal/models"
	errors2 "DriverService/internal/trip_errors"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	trips *mongo.Collection
}

func NewRepository(trips *mongo.Collection) (*Repository, error) {
	return &Repository{
		trips: trips,
	}, nil
}

func (repo *Repository) GetAllTrips() ([]models.Trip, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := repo.trips.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var trips []models.Trip
	err = cur.All(ctx, &trips)
	return trips, err
}

func (repo *Repository) Get(id string) (*models.Trip, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var trip models.Trip
	err := repo.trips.FindOne(ctx, bson.M{"id": id}).Decode(&trip)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors2.NotFoundError{Key: id}
		}
		return nil, err
	}
	return &trip, nil
}

func (repo *Repository) Add(value models.Trip) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.trips.InsertOne(ctx, value)
	return err
}

func (repo *Repository) ChangeTripStatus(id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := repo.trips.FindOneAndUpdate(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors2.NotFoundError{Key: id}
		}
		return err
	}
	return nil
}
