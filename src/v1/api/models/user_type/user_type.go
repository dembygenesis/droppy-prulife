package user_type

import (
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (u *UserType) GetNameById() (string, error) {
	var userType string

	sql := `
		SELECT 
		    name
		FROM user_type
		WHERE 1 = 1 
			AND id = ?
	`

	err := database.DBInstancePublic.Get(&userType, sql, u.ID)

	return userType, err
}

func (u *UserType) GetAll() ([]UserType, error) {
	var userTypes []UserType
	sql := `
		SELECT 
			id,
		    name
		FROM user_type
		WHERE 1 = 1 
	`

	err := database.DBInstancePublic.Select(&userTypes, sql)

	return userTypes, err
}

func (u *UserType) ValidID() (bool, error) {
	hasId := false
	sql := `
		SELECT 
			IF(COUNT(id) > 0, true, false) AS has_id 
		FROM user_type
		WHERE 1 = 1
			AND id = ? 
	`

	err := database.DBInstancePublic.Get(&hasId, sql, u.ID)

	return hasId, err
}