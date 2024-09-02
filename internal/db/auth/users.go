package auth

import (
	"crypto/sha256"
	"encoding/hex"

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
	hashedKey := sha256.New()
	hashedKey.Write([]byte(apiKey))

	hexHashedKey := hex.EncodeToString(hashedKey.Sum(nil))

	var user db.User
	err := dbClient.Preload(clause.Associations).Where(&db.User{APIKey: hexHashedKey}).First(&user).Error
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
	hashedKey := sha256.New()
	hashedKey.Write([]byte(newApiKey))
	hexHashedKey := hex.EncodeToString(hashedKey.Sum(nil))

	res := dbClient.
		Where("id", userID).
		Updates(db.User{
			APIKey: hexHashedKey,
		})
	if res.Error != nil {
	  return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return newApiKey, nil
}