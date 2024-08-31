package auth

import (
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

func CreateOrUpdateUser(dbClient *gorm.DB, userData db.User) (*db.User, error) {
	err := dbClient.FirstOrCreate(&userData).Error
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

// TODO Rotate API Key function