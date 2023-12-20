package mongo

import (
	"DriverService/internal/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotFoundError struct {
	Key string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Key)
}

type DuplicateKeyError struct {
	Key string
}

func (e DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate key: %s", e.Key)
}

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
}

func NewRepository(uri string, dbName string, collName string) (*Repository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collName)

	return &Repository{
		client:     client,
		collection: collection,
		ctx:        context.Background(),
	}, nil
}

func (repo *Repository) GetAllTrips() ([]models.Trip, error) {
	cur, err := repo.collection.Find(repo.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var trips []models.Trip
	err = cur.All(repo.ctx, &trips)
	if err != nil {
		return nil, err
	}
	return trips, nil
}

func (repo *Repository) Get(id string) (models.Trip, error) {
	var trip models.Trip
	err := repo.collection.FindOne(repo.ctx, bson.M{"id": id}).Decode(&trip)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return trip, NotFoundError{id}
		}
		return trip, err
	}
	return trip, nil
}

func (repo *Repository) Add(value models.Trip) error {
	_, err := repo.collection.InsertOne(repo.ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) ChangeTripStatus(id string, status string) error {
	result := repo.collection.FindOneAndUpdate(repo.ctx,
		bson.M{"id": id},
		bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}}}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return NotFoundError{id}
	}
	return result.Err()
}
