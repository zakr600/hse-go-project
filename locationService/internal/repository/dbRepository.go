package repository

import (
	"database/sql"
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/config"
	"github.com/lib/pq"
	"log"
)

type DbRepository struct {
	db *sql.DB
}

func CreateDbRepository(config config.DbConfig) *DbRepository {
	rawConnStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Username, config.Password, config.DbPackageName, config.DbPort, config.Dbname)
	connStr, err := pq.ParseURL(rawConnStr)
	db, err := sql.Open(config.DbDriverName, connStr)
	if err != nil {
		panic(err)
	}
	return &DbRepository{
		db: db,
	}
}

func (r *DbRepository) SetDriverLocation(driverId string, location Location) {
	r.db.Exec("UPDATE locations SET lat=$1, lng=$2 WHERE driver_id=$3", location.Lat, location.Lng, driverId)
	r.db.Exec("INSERT INTO locations (driver_id, lat, lng) SELECT $3, $1, $2 WHERE NOT EXISTS (SELECT 1 FROM locations WHERE driver_id=$3)", location.Lat, location.Lng, driverId)
}

func (r *DbRepository) GetAllDrivers() []Driver {
	rows, err := r.db.Query("SELECT driver_id, lat, lng FROM locations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	drivers := make([]Driver, 0)
	for rows.Next() {
		var driver Driver
		err := rows.Scan(&driver.DriverId, &driver.Location.Lat, &driver.Location.Lng)
		if err != nil {
			log.Fatal(err)
		}
		drivers = append(drivers, driver)
	}

	return drivers
}
