package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Migration struct {
	client *mongo.Client
	db     *mongo.Database
	log    *zap.Logger
}

func NewMigration(client *mongo.Client, db *mongo.Database, log *zap.Logger) *Migration {
	return &Migration{
		client: client,
		db:     db,
		log:    log,
	}
}
func (m *Migration) Run(path string) error {
	m.log.Info("Running migration")
	dbDriver, err := mongodb.WithInstance(
		m.client,
		&mongodb.Config{DatabaseName: "driver_service"},
	)
	if err != nil {
		return err
	}

	fileSrc, err := (&file.File{}).Open(path)
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance("file", fileSrc, "driver_service", dbDriver)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil {
		return err
	}

	return nil
}
