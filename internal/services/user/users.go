package user

import (
	"github.com/graytonio/flagops-data-store/internal/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Handles managing user data including permissions and authentication source
type UserDataService struct {
	DBClient *gorm.DB
}

// TODO Add pagination
func (ud *UserDataService) GetUsers() ([]db.User, error) {
	users := []db.User{}

	err := ud.DBClient.Preload(clause.Associations).Find(&users).Error
	if err != nil {
	  return nil, err
	}

	return users, nil
}

func (ud *UserDataService) GetUserByID(id uint) (*db.User, error) {
	user := db.User{}
	res := ud.DBClient.Preload(clause.Associations).First(&user, id)
	if res.Error != nil {
	  return nil, res.Error
	}

	return &user, nil
}


func (ud *UserDataService) UpsertUser(userData db.User) (*db.User, error) {
	err := ud.DBClient.Preload(clause.Associations).Clauses(clause.OnConflict{
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
