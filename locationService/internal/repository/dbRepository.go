package repository

import (
	"database/sql"
	"fmt"
)

type dbRepository struct {
	db *sql.DB
}

type DbConfig struct {
	username   string
	password   string
	dbname     string
	driverName string
}

func createDbRepository(config DbConfig) *dbRepository {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.username, config.password, config.dbname)
	db, err := sql.Open(config.driverName, connStr)
	if err != nil {
		panic(err)
	}
	return &dbRepository{
		db: db,
	}
}

func (r *dbRepository) SetDriverLocation(driverId string, location Location) {

}
