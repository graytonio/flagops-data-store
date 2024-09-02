package auth

import (
	"github.com/graytonio/flagops-data-storage/internal/db"
	"gorm.io/gorm"
)

func GetPermissions(dbClient *gorm.DB) ([]db.Permission, error) {
	permissions := []db.Permission{}

	err := dbClient.Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
