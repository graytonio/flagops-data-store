package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO Test db operations

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
			OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"display_name",
				}),
			}).
			Create(&BootstrapPermissions).Error
}