package orders

import (
	"database/sql"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (o *Order) GetAll() {

}

func (o *Order) GetDetails(orderId int, userId int, userType string) (*[]ResponseOrdersDisplay, error) {

	var responseOrdersDisplay []ResponseOrdersDisplay
	var err error

	sql := ""

	if userType == "Seller" {
		sql = `
			 SELECT
			  o.id,
			  o.amount,
			  IF(DATE_FORMAT(o.created_date, '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(o.created_date, '%Y-%m-%d %h:%i %p')) AS date_created,
			  IF(CONCAT(u_seller.lastname, ', ', u_seller.firstname) IS NULL, "", CONCAT(u_seller.lastname, ', ', u_seller.firstname)) AS seller,
			  p.name AS product,
			  od.price_per_item,
			  od.quantity,
			  od.total_price,
			  o.is_active,
    	      o.amount AS order_total_price,
    	      r.name AS region,
    	      CONCAT(u_dropshipper.lastname, ', ', u_dropshipper.firstname) AS dropshipper
			FROM
			  `+ "`order`" + `o
			  INNER JOIN `+ "`order_detail`" + `od
				ON 1 = 1
				AND o.id = od.order_id
			  INNER JOIN product p
				ON 1 = 1
				AND od.product_id = p.id
    	      INNER JOIN region r
				ON 1 = 1
				AND o.region_id = r.id
			  INNER JOIN user u_seller
				ON 1 = 1
				  AND o.seller_id = u_seller.id
    	      INNER JOIN user u_dropshipper 
				ON 1 = 1
    	          AND o.dropshipper_id = u_dropshipper.id
    	    WHERE 1 = 1
				AND o.id = ? 
    	        AND o.seller_id = ?
				AND o.is_active = 1
		`

		err = database.DBInstancePublic.Select(&responseOrdersDisplay, sql, orderId, userId)
	}

	if userType == "Dropshipper" {
		sql = `
			 SELECT
			  o.id,
			  o.amount,
			  IF(DATE_FORMAT(o.created_date, '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(o.created_date, '%Y-%m-%d %h:%i %p')) AS date_created,
			  IF(CONCAT(u_seller.lastname, ', ', u_seller.firstname) IS NULL, "", CONCAT(u_seller.lastname, ', ', u_seller.firstname)) AS seller,
			  p.name AS product,
			  od.price_per_item,
			  od.quantity,
			  od.total_price,
			  o.is_active,
    	      o.amount AS order_total_price,
    	      r.name AS region,
    	      CONCAT(u_dropshipper.lastname, ', ', u_dropshipper.firstname) AS dropshipper
			FROM
			  `+ "`order`" + `o
			  INNER JOIN `+ "`order_detail`" + `od
				ON 1 = 1
				AND o.id = od.order_id
			  INNER JOIN product p
				ON 1 = 1
				AND od.product_id = p.id
    	      INNER JOIN region r
				ON 1 = 1
				AND o.region_id = r.id
			  INNER JOIN user u_seller
				ON 1 = 1
				  AND o.seller_id = u_seller.id
    	      INNER JOIN user u_dropshipper 
				ON 1 = 1
    	          AND o.dropshipper_id = u_dropshipper.id
    	    WHERE 1 = 1
				AND o.id = ? 
    	        AND o.dropshipper_id = ?
				AND o.is_active = 1
		`

		err = database.DBInstancePublic.Select(&responseOrdersDisplay, sql, orderId, userId)
	}

	if userType == "Admin" {
		sql = `
			 SELECT
			  o.id,
			  o.amount,
			  IF(DATE_FORMAT(o.created_date, '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(o.created_date, '%Y-%m-%d %h:%i %p')) AS date_created,
			  IF(CONCAT(u_seller.lastname, ', ', u_seller.firstname) IS NULL, "", CONCAT(u_seller.lastname, ', ', u_seller.firstname)) AS seller,
			  p.name AS product,
			  od.price_per_item,
			  od.quantity,
			  od.total_price,
			  o.is_active,
    	      o.amount AS order_total_price,
    	      r.name AS region,
    	      CONCAT(u_dropshipper.lastname, ', ', u_dropshipper.firstname) AS dropshipper
			FROM
			  `+ "`order`" + `o
			  INNER JOIN `+ "`order_detail`" + `od
				ON 1 = 1
				AND o.id = od.order_id
			  INNER JOIN product p
				ON 1 = 1
				AND od.product_id = p.id
    	      INNER JOIN region r
				ON 1 = 1
				AND o.region_id = r.id
			  INNER JOIN user u_seller
				ON 1 = 1
				  AND o.seller_id = u_seller.id
    	      INNER JOIN user u_dropshipper 
				ON 1 = 1
    	          AND o.dropshipper_id = u_dropshipper.id
    	    WHERE 1 = 1
				AND o.id = ? 
				AND o.is_active = 1
		`

		err = database.DBInstancePublic.Select(&responseOrdersDisplay, sql, orderId)
	}

	return &responseOrdersDisplay, err
}

func (o *Order) Create(p *ParamsOrder) (*sql.Result, error) {

	sql := `
		CALL add_order(
		  ?,
		  ?,
		  ?
		);
	`

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		p.AdminId,
		p.UserId,
		p.OrderDetails,
	)

	return &sqlResult, err
}

func (o *Order) Update(p *ParamsOrderUpdate) (*sql.Result, error) {

	fmt.Println("p", p)

	sql := `
		CALL update_order(
		  ?,
		  ?,
		  ?
		);
	`

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		p.AdminId,
		p.OrderId,
		p.Description,
	)

	return &sqlResult, err
}