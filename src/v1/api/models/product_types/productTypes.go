package product_types

import (
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)


func (p *ProductType) GetAll() (*[]ProductType, error) {
	var productTypes []ProductType

	sql := `
		SELECT 
			id,
		    name
		FROM product_type
		WHERE 1 = 1
	`

	err := database.DBInstancePublic.Select(&productTypes, sql)

	return &productTypes, err
}

func (p *ProductType) ValidId() (bool, error) {
	var isTrue bool

	sql := `
		SELECT 
			IF(COUNT(id) > 0, true, false) AS is_true 
		FROM product_type
		WHERE 1 = 1
			AND id = ?
	`

	err := database.DBInstancePublic.Get(&isTrue, sql, p.ID)

	return isTrue, err
}

func (p *ProductType) UniqueName() (bool, error) {
	var isTrue bool

	sql := `
		SELECT 
			IF(COUNT(id) = 0, true, false) AS is_true 
		FROM product_type
		WHERE 1 = 1
			AND name = ?
	`

	err := database.DBInstancePublic.Get(&isTrue, sql, p.Name)

	fmt.Println("isTrue", isTrue)

	return isTrue, err
}