package products

import (
	"database/sql"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (p *Product) GetSellerList(userId int) (*[]ResponseSellerList, error) {
	var responseSellerList []ResponseSellerList

	sql := `
		SELECT
		  u.id,
		  CONCAT(u.firstname, ', ', u.lastname) AS name
		FROM
		  inventory i
		  INNER JOIN user u 
			ON 1 = 1
			  AND i.seller_id = u.id
		WHERE 1 = 1
		  AND i.dropshipper_id = ?
		  AND u.is_active = 1
	`

	err := database.DBInstancePublic.Select(&responseSellerList, sql, userId)

	return &responseSellerList, err
}

func (p *Product) GetInventory(userId int, userType string) ([]ResponseInventoryList, error) {
	var responseInventoryList []ResponseInventoryList

	sql := ""

	if userType == "Seller" {
		sql = `
			SELECT
			  p.id,
			  p.name,
			  pt.name AS category,
			  SUM(i.quantity) AS remaining,
			  r.name AS region
			FROM
			  inventory i
			  INNER JOIN product p
				ON 1 = 1
				AND i.product_id = p.id
			  INNER JOIN product_type pt
				ON 1 = 1
				AND p.product_type_id = pt.id
			  INNER JOIN region r 
				ON 1 = 1
				AND i.region_id = r.id
			WHERE 1 = 1
			  AND i.is_active = 1 
			  AND p.is_active = 1 
			  AND i.seller_id = ?
			GROUP BY p.id, r.id
		`
	}

	if userType == "Dropshipper" {
		sql = `
			SELECT
			  p.id,
			  p.name,
			  pt.name AS category,
			  SUM(i.quantity) AS remaining,
			  r.name AS region
			FROM
			  inventory i
			  INNER JOIN product p
				ON 1 = 1
				AND i.product_id = p.id
			  INNER JOIN product_type pt
				ON 1 = 1
				AND p.product_type_id = pt.id
			  INNER JOIN region r 
				ON 1 = 1
				AND i.region_id = r.id
			WHERE 1 = 1
			  AND i.is_active = 1 
			  AND p.is_active = 1 
			  AND i.dropshipper_id = ?
			GROUP BY p.id, r.id
		`
	}

	err := database.DBInstancePublic.Select(&responseInventoryList, sql, userId)

	if err != nil {
		return responseInventoryList, err
	}
	return responseInventoryList, err
}

func (p *Product) ValidId() (bool, error) {
	var isTrue bool

	sql := `
		SELECT 
			IF(COUNT(id) > 0, true, false) AS is_true 
		FROM product
		WHERE 1 = 1
			AND id = ?
	`

	err := database.DBInstancePublic.Get(&isTrue, sql, p.ID)

	return isTrue, err
}

func (p *Product) UpdateUrl() (sql.Result, error) {
	sql := `
		UPDATE product 
		SET url = ?
		WHERE id = ?
	`

	res, err := database.DBInstancePublic.Exec(sql, p.Url, p.ID)

	fmt.Println("yawa", res, err)

	return res, err
}

func (p *Product) Update() (sql.Result, error) {
	sql := `
		UPDATE product 
		SET name = ?,
			product_type_id = ?
		WHERE id = ?
	`

	res, err := database.DBInstancePublic.Exec(sql, p.Name, p.ProductTypeId, p.ID)

	return res, err
}

func (p *Product) GetAll(filter string) ([]ResponseProductList, error) {
	var productList []ResponseProductList

	exclusions := "1 = 1"

	if filter == "package" {
		exclusions = "p.name != (SELECT `value` FROM sysparam WHERE `key` = 'PRODUCT_EXCLUSIONS_PACKAGE')"
	}

	if filter == "parcel" {
		exclusions = "p.name != (SELECT `value` FROM sysparam WHERE `key` = 'PRODUCT_EXCLUSIONS_PARCEL')"
	}

	if filter == "dropship" {
		exclusions = "p.name != (SELECT `value` FROM sysparam WHERE `key` = 'PRODUCT_EXCLUSIONS_DROPSHIP')"
	}

	sql := `
		SELECT 
		  p.id,
		  IF(p.url IS NULL, "", p.url) AS url,
		  IF(p.name IS NULL, "", p.name) AS name,
		  IF(pt.name IS NULL, "", pt.name) AS category,
		  IF(pt.id IS NULL, 0, pt.id) AS category_id, 
		  p.price_per_item,
		  p.price_per_item_dropshipper
		FROM
		  product p 
		  INNER JOIN product_type pt 
			ON 1 = 1 
			AND p.product_type_id = pt.id 
		WHERE 1 = 1 
		  AND p.is_active = 1
  		  AND ` + exclusions + `
	`

	fmt.Println(sql)

	err := database.DBInstancePublic.Select(&productList, sql)

	if err != nil {
		return productList, err
	}
	return productList, err
}

func (p *Product) GetOne() (*[]ResponseProductSpecific, error) {
	var responseProductSpecific []ResponseProductSpecific

	sql := `
		SELECT 
		  p.id,
		  p.name,
		  p.url,
		  pt.name AS category,
		  pt.id AS product_type_id
		FROM
		  product p 
		  INNER JOIN product_type pt 
			ON 1 = 1 
			AND p.product_type_id = pt.id 
		WHERE 1 = 1 
		  AND p.is_active = 1 
		  AND p.id = ?
	`

	err := database.DBInstancePublic.Select(
		&responseProductSpecific,
		sql,
		p.ID,
		)

	return &responseProductSpecific, err
}

func (p *Product) Create() (sql.Result, error) {

	sql := `
		INSERT INTO product (
		  name,
		  product_type_id,
		  created_by,
		  created_date,
		  url
		)
		VALUES
		  (
			?,
			?,
			?,
			NOW(),
			1
		  );
	`

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		p.Name,
		p.ProductTypeId,
		p.CreatedBy,
	)

	fmt.Println(sqlResult, err)

	return sqlResult, err
}

func (p *Product) UniqueName() (bool, error) {
	var isTrue bool

	sql := `
		SELECT 
			IF(COUNT(id) = 0, true, false) AS is_true 
		FROM product
		WHERE 1 = 1
			AND name = ?
	`

	err := database.DBInstancePublic.Get(&isTrue, sql, p.Name)

	fmt.Println("isTrue", isTrue)

	return isTrue, err
}

func (p *Product) UniqueNameExceptOwn() (bool, error) {
	var isTrue bool

	sql := `
		SELECT 
			IF(COUNT(id) = 0, true, false) AS is_true 
		FROM product
		WHERE 1 = 1
			AND name = ?
			AND id != ?
	`

	err := database.DBInstancePublic.Get(&isTrue, sql, p.Name, p.ID)

	fmt.Println("isTrue", isTrue)

	return isTrue, err
}

func (p *Product) Delete() (sql.Result, error) {

	sql := `
		UPDATE product
		SET is_active = 0
		WHERE id = ?
	`

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		p.ID,
	)

	fmt.Println(sqlResult, err)

	return sqlResult, err
}