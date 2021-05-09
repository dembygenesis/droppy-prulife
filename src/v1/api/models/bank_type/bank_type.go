package bank_type

import (
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (u *BankType) GetAll() ([]BankType, error) {
	var bankTypes []BankType

	sql := `
		SELECT 
			id, 
		    name
		FROM bank_type
		WHERE 1 = 1 
	`

	err := database.DBInstancePublic.Select(&bankTypes, sql)

	return bankTypes, err
}

func (u *BankType) ValidID() (bool, error) {
	hasId := false
	sql := `
		SELECT 
			IF(COUNT(id) > 0, true, false) AS has_id 
		FROM bank_type
		WHERE 1 = 1
			AND id = ? 
	`

	err := database.DBInstancePublic.Get(&hasId, sql, u.ID)

	return hasId, err
}