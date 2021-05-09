package region

import (
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (u *Region) GetAll() ([]Region, error) {
	var regions []Region

	fmt.Println("Hello")

	sql := `
		SELECT 
			id, 
		    name
		FROM ` + "`region`" + `
		WHERE 1 = 1 
	`

	err := database.DBInstancePublic.Select(&regions, sql)

	return regions, err
}

func (u *Region) ValidID() (bool, error) {
	hasId := false
	sql := `
		SELECT 
			IF(COUNT(id) > 0, true, false) AS has_id 
		FROM region
		WHERE 1 = 1
			AND id = ? 
	`

	err := database.DBInstancePublic.Get(&hasId, sql, u.ID)

	fmt.Println("hasId", hasId)
	fmt.Println("u.ID", u.ID)

	return hasId, err
}