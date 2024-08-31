package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Creates a new db client and runs auto migration on it
func GetDBClient(dsn string) (*gorm.DB, error) {
	dbClient, err :=  gorm.Open(postgres.Open(dsn))
	if err != nil {
	  return nil, err
	}

	err = dbClient.AutoMigrate(&User{}, &Permission{})
	if err != nil {
	  return nil, err
	}

	err = bootstrapPermissions(dbClient)
	if err != nil {
	  return nil, err
	}

	return dbClient, nil
}

func bootstrapPermissions(db *gorm.DB) error {
	return db.Clauses(clause.
			OnConflict{DoNothing: true}). // TODO On conflict update display name
			Create(&BootstrapPermissions).Error
}