package auth

import (
	"github.com/google/uuid"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUsers(dbClient *gorm.DB) ([]db.User, error) {
	users := []db.User{}

	err := dbClient.Preload(clause.Associations).Find(&users).Error
	if err != nil {
	  return nil, err
	}

	return users, nil
}

func GetUserByID(dbClient *gorm.DB, id uint) (*db.User, error) {
	user := db.User{}
	res := dbClient.Preload(clause.Associations).First(&user, id)
	if res.Error != nil {
	  return nil, res.Error
	}

	return &user, nil
}

func GetUserByAPIKey(dbClient *gorm.DB, apiKey string) (*db.User, error) {
	var user db.User
	err := dbClient.Preload(clause.Associations).Where(&db.User{APIKey: apiKey}).First(&user).Error // TODO Hash api key
	if err != nil {
	  return nil, err
	}

	return &user, err
}

func UpsertUser(dbClient *gorm.DB, userData db.User) (*db.User, error) {
	err := dbClient.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"email",
			"username",
		}),
	}).FirstOrCreate(&userData).Error
	if err != nil {
	  return nil, err
	}

	return &userData, nil
}

func AddUserPermissions(dbClient *gorm.DB, userID uint, permissionIDs []string) error {
	user, err := GetUserByID(dbClient, userID)
	if err != nil {
	  return err
	}


	permissions := []db.Permission{}
	for _, p := range permissionIDs{
		permissions = append(permissions, db.Permission{ID: p})
	}

	err = dbClient.Model(&user).Association("Permissions").Append(permissions)
	if err != nil {
	  return err
	}

	return nil
}

func RemoveUserPermissions(dbClient *gorm.DB, userID uint, permissionIDs []string) error {
	user, err := GetUserByID(dbClient, userID)
	if err != nil {
	  return err
	}


	permissions := []db.Permission{}
	for _, p := range permissionIDs{
		permissions = append(permissions, db.Permission{ID: p})
	}

	err = dbClient.Model(&user).Association("Permissions").Delete(permissions)
	if err != nil {
	  return err
	}

	return nil
}

func RotateUserAPIKey(dbClient *gorm.DB, userID uint) (string, error) {
	newApiKey := uuid.NewString()

	err := dbClient.
		Where("id", userID).
		Updates(db.User{
			APIKey: newApiKey, // TODO Hash API Key
		}).Error
	if err != nil {
	  return "", err
	}
	return newApiKey, nil
}