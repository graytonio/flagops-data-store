package user

import (
	"github.com/graytonio/flagops-data-store/internal/db"
)

// TODO Add pagination
func (ud *UserDataService) GetPermissions() ([]db.Permission, error) {
	permissions := []db.Permission{}

	err := ud.DBClient.Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (ud *UserDataService) AddUserPermissions(userID uint, permissionIDs []string) error {
	user, err := ud.GetUserByID(userID)
	if err != nil {
	  return err
	}


	permissions := []db.Permission{}
	for _, p := range permissionIDs{
		permissions = append(permissions, db.Permission{ID: p})
	}

	err = ud.DBClient.Model(&user).Association("Permissions").Append(permissions)
	if err != nil {
	  return err
	}

	return nil
}

func (ud *UserDataService) RemoveUserPermissions(userID uint, permissionIDs []string) error {
	user, err := ud.GetUserByID(userID)
	if err != nil {
	  return err
	}


	permissions := []db.Permission{}
	for _, p := range permissionIDs{
		permissions = append(permissions, db.Permission{ID: p})
	}

	err = ud.DBClient.Model(&user).Association("Permissions").Delete(permissions)
	if err != nil {
	  return err
	}

	return nil
}
