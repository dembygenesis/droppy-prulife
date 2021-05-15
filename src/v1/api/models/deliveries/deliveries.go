package deliveries

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
	UserModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/users"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"strconv"
)

func (d *Delivery) UpdateDelivery(p ParamsUpdateDelivery) (sql.Result, error) {
	sql := `CALL update_delivery(?, ?, ?, ?, ?);`

	res, err := database.DBInstancePublic.Exec(
		sql,
		d.UserId,
		p.DeliveryId,
		p.DeliveryStatus,
		p.TrackingNumber,
		p.VoidOrRejectReason,
	)

	return res, err
}

func (d *Delivery) GetDetails() {
	// var responseDeliveryDetails ResponseDeliveryDetails
}

func (d *Delivery) CreateParcel(p ParamsCreateParcel) (sql.Result, error) {
	sql := `
		CALL add_delivery(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	fmt.Println(p)

	res, err := database.DBInstancePublic.Exec(sql,
		p.DeliveryOption,
		p.SellerId,
		p.DropshipperId,
		p.Name,
		p.ContactNumber,
		p.Address,
		p.Note,
		p.Region,
		p.ServiceFee,
		p.DeclaredAmount,
		p.DeliveryDetails,
	)

	return res, err
}

func (d *Delivery) CreatePackage(p ParamsCreateOrder) (sql.Result, error) {

	// These will change based on the emails LATER

	sql := `
		CALL add_order(?, ?, ?, ?)
	`

	res, err := database.DBInstancePublic.Exec(sql,
		p.SellerId,
		p.DropshipperId,
		p.OrderDetails,
		p.Region,
	)

	return res, err
}

func (d *Delivery) GetCoinTransactions(userId int, userType string, page int, rows int) (*[]ResponseCoinTransactions, response_builder.Pagination, error) {
	var container []ResponseCoinTransactions
	var pagination response_builder.Pagination

	sql := `
		SELECT 
		  date_created,		
		  ` + "`type`" + `,
		  amount,
		  reference_number,
		  bank_type,
		  source,
		  recipient,
		  tran_num,
		  is_active
		FROM (
		  (
		    SELECT
		      IF(DATE_FORMAT(CONVERT_TZ(t.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(t.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
		      IF (t.money_in = 1, 'Money In', 'Money Out') AS type, 
			  t.created_date,
		      t.amount,
		      t.reference_number,
		      bt.name AS bank_type,
			  'Cash In' AS source,
			  'N/A' AS recipient,
			  t.id AS tran_num,
			  t.is_active
		    FROM
		      ` + "`transaction`" + ` t
		      INNER JOIN bank_type bt
		        ON 1 =1
		          AND t.bank_type_id = bt.id
		      WHERE 1 = 1
		        AND t.user_allotted_id = ?
		        AND t.is_active = 1
			ORDER BY t.created_date DESC
		  )
		  UNION ALL 
		  (
		    SELECT
		      IF(DATE_FORMAT(CONVERT_TZ(ct.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(ct.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
		      IF(ct.type = 'C', 'Coins In', 'Coins Out') AS type,
			  ct.created_date,
		      ct.amount * -1 AS amount,
		      'N/A' AS reference_number,
		      'N/A' AS bank_type,
			  CASE
			    WHEN ct.withdrawal_id IS NOT NULL THEN 'Withdrawal'	
			    WHEN ct.delivery_id IS NOT NULL THEN 'Delivery'	
			    WHEN ct.order_id IS NOT NULL THEN 'Order'	
			    WHEN ct.coin_transaction_id IS NOT NULL THEN 'Coins added from Cash In'
              END AS source,
			  IF (d.id IS NULL, 'N/A', d.name) AS recipient,
			  CASE
			    WHEN ct.withdrawal_id IS NOT NULL THEN ct.withdrawal_id	
			    WHEN ct.delivery_id IS NOT NULL THEN ct.delivery_id	
			    WHEN ct.order_id IS NOT NULL THEN ct.order_id	
			    WHEN ct.coin_transaction_id IS NOT NULL THEN ct.coin_transaction_id
              END AS tran_num,
			  ct.is_active
		    FROM
		      coin_transaction ct
		      LEFT JOIN delivery d 
			    ON 1 = 1
					AND ct.delivery_id = d.id 
		      WHERE 1 = 1
		        AND ct.user_id = ?
		        AND ct.is_active = 1
			ORDER BY ct.created_date DESC
		  )
		) AS a
		ORDER BY created_date DESC
	`

	var err error
	var count int

	paginate := func () (*[]ResponseCoinTransactions, response_builder.Pagination, error) {
		count, err = database.GetQueryCount(sql, userId, userId)

		if err != nil {
			return &container, pagination, err
		}

		if count == 0 {
			return &container, pagination, nil
		}

		sql, pages, rowsPerPage, offset, page, totalCount, resultCount := database.GetPaginationDetails(
			sql,
			count,
			page,
			rows,
			1000,
		)

		pagination.SetData(rowsPerPage, offset, pages, rows, page, totalCount, resultCount)

		err := database.DBInstancePublic.Select(
			&container,
			sql,
			userId,
			userId,
		)

		if err != nil {
			return &container, pagination, err
		}

		return &container, pagination, nil
	}

	res, pagination, err := paginate()

	return res, pagination, err
}

func (d *Delivery) GetCoinTransactions2(userId int, userType string, page int, rows int) (*[]ResponseCoinTransactions, response_builder.Pagination, error) {
	var container []ResponseCoinTransactions
	var pagination response_builder.Pagination

	sql := `
		SELECT 
		  date_created,		
		  ` + "`type`" + `,
		  amount,
		  reference_number,
		  bank_type,
		  source,
		  recipient,
		  tran_num,
		  void_or_reject_reason,
		  CASE
            WHEN type IN ('Money In', 'Money Out')
            THEN @running_total
            ELSE @running_total := @running_total + amount
          END AS running_balance,
          created_date,
          source_type,
		  is_active
		FROM (
		  (
		    SELECT
		      IF(DATE_FORMAT(CONVERT_TZ(t.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(t.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
		      IF (t.money_in = 1, 'Money In', 'Money Out') AS type, 
			  t.created_date,
		      t.amount,
		      t.reference_number,
		      bt.name AS bank_type,
			  'Cash In' AS source,
			  'N/A' AS recipient,
			  t.id AS tran_num,
              'N/A' AS void_or_reject_reason,
              'N/A' AS source_type,
			  t.is_active
		    FROM
		      ` + "`transaction`" + ` t
		      INNER JOIN bank_type bt
		        ON 1 =1
		          AND t.bank_type_id = bt.id
		      WHERE 1 = 1
		        AND t.user_allotted_id = ?
		        AND t.is_active = 1
			ORDER BY t.created_date DESC
		  )
		  UNION ALL 
		  (
		    SELECT
		      IF(DATE_FORMAT(CONVERT_TZ(ct.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(ct.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
		      IF(ct.type = 'C', 'Coins In', 'Coins Out') AS type,
			  ct.created_date,
		      ct.amount * -1 AS amount,
		      'N/A' AS reference_number,
		      'N/A' AS bank_type,
			  CASE
			    WHEN ct.withdrawal_id IS NOT NULL THEN 'Withdrawal'	
			    WHEN ct.delivery_id IS NOT NULL THEN 'Delivery'	
			    WHEN ct.order_id IS NOT NULL THEN 'Order'	
			    WHEN ct.coin_transaction_id IS NOT NULL THEN 'Coins added from Cash In'
              END AS source,
			  IF (d.id IS NULL, 'N/A', d.name) AS recipient,
			  CASE
			    WHEN ct.withdrawal_id IS NOT NULL THEN ct.withdrawal_id	
			    WHEN ct.delivery_id IS NOT NULL THEN ct.delivery_id	
			    WHEN ct.order_id IS NOT NULL THEN ct.order_id	
			    WHEN ct.coin_transaction_id IS NOT NULL THEN ct.coin_transaction_id
              END AS tran_num,
			  CASE
			    WHEN ct.withdrawal_id IS NOT NULL THEN 'N/A'
			    WHEN ct.delivery_id IS NOT NULL THEN IF (d.void_or_reject_reason IS NULL, '', d.void_or_reject_reason)
			    WHEN ct.order_id IS NOT NULL THEN 'N/A'	
			    WHEN ct.coin_transaction_id IS NOT NULL THEN 'N/A'
              END AS void_or_reject_reason,
			  CASE
			    WHEN ct.delivery_id IS NOT NULL THEN ds.name
 			    ELSE 'N/A'
              END AS source_type,
			  ct.is_active
		    FROM
		      coin_transaction ct
		      LEFT JOIN delivery d 
			    ON 1 = 1
					AND ct.delivery_id = d.id
			  LEFT JOIN delivery_status ds 
			    ON 1 = 1
					AND d.delivery_status_id = ds.id
		      WHERE 1 = 1
		        AND ct.user_id = ?
		        AND ct.is_active = 1
			ORDER BY ct.created_date DESC
		  )
		) AS a
		ORDER BY created_date ASC
	`

	var err error
	var count int

	paginate := func () (*[]ResponseCoinTransactions, response_builder.Pagination, error) {
		count, err = database.GetQueryCount(sql, userId, userId)

		if err != nil {
			return &container, pagination, err
		}

		if count == 0 {
			return &container, pagination, nil
		}

		sql, pages, rowsPerPage, offset, page, totalCount, resultCount := database.GetPaginationDetails(
			sql,
			count,
			page,
			rows,
			1000,
		)

		pagination.SetData(rowsPerPage, offset, pages, rows, page, totalCount, resultCount)

		/**
			Replace with transactions here
		 */

		sql = `
			SELECT * FROM (
				` + sql + `    
			) AS a
			ORDER BY created_date DESC
		`

		tx := database.DBInstancePublic.MustBegin()

		_, err := tx.Exec("SET @running_total = 0")

		if err != nil {
			return &container, pagination, errors.New("SET @running_total = 0 problem")
		}

		err = tx.Select(&container, sql, userId, userId)

		return &container, pagination, err
	}

	res, pagination, err := paginate()

	return res, pagination, err
}

// Should I add a new filter object? And validate? LOL... Maybe.
func (d *Delivery) GetTransactions(
	userId int,
	userType string,
	deliveryStatus string,
	search SearchTransactionFilter,
	page int,
	rows int,
) (*[]ResponseTransactions, response_builder.Pagination, error) {
	// var container []ResponseTransactions

	container := make([]ResponseTransactions, 0)

	var pagination response_builder.Pagination

	var sqlUserFilter string
	var sqlTranNumFilter string

	// Handle admin search filters
	userFilter := search.UserId
	tranNumFilter := search.TranNum

	sql := ""

	if userFilter != 0 {
		sqlUserFilter = "u.id = " + strconv.Itoa(userFilter)
	} else {
		sqlUserFilter = "1 = 1"
	}

	if tranNumFilter != 0 {
		sqlTranNumFilter = "d.id = " + strconv.Itoa(tranNumFilter)
	} else {
		sqlTranNumFilter = "1 = 1"
	}

	// ADMINS
	if userType == "Admin" {
		sql = `
			SELECT 
			  date_created,
			  recipient,
			  transaction_number,
			  CAST(amount AS DECIMAL(65,2)) AS amount,
			  tracking_number,
			  ` + "`type`" + `,  
			  ` + "`status`" + `,
			  region,
			  seller,
			  dropshipper,
			  delivery_payment_method,
			  d.contact_number,
			  d.item_description
			FROM (
				(
					SELECT 
                      d.created_date,
					  IF(DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
					  d.name AS recipient,
					  d.id AS transaction_number,
					  d.declared_amount AS amount,
					  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
					  do.name AS type,
					  ds.name AS status,
 					  CONCAT(u_dropshipper.lastname, ', ', u_dropshipper.firstname) AS dropshipper,
 					  CONCAT(u_seller.lastname, ', ', u_seller.firstname) AS seller,
					  r.name AS region,
 					  (SELECT COUNT(id) FROM delivery_detail WHERE delivery_id = d.id) AS items,
					  dpm.name AS delivery_payment_method
					FROM delivery d
					INNER JOIN delivery_status ds 
					  ON 1 = 1
						AND d.delivery_status_id = ds.id 
					INNER JOIN delivery_option do
					  ON 1 = 1 
						AND d.delivery_option_id = do.id
					INNER JOIN ` + "`user`" + ` u_dropshipper
 					  ON 1 = 1
						AND d.dropshipper_id = u_dropshipper.id
					INNER JOIN ` + "`user`" + ` u_seller
 					  ON 1 = 1
						AND d.seller_id = u_seller.id
				    INNER JOIN region r
 					  ON 1 = 1
						AND d.region_id = r.id
					INNER JOIN delivery_payment_method dpm
 					  ON 1 = 1
						AND d.delivery_payment_method_id = dpm.id
					WHERE 1 = 1
					 AND (d.is_active = 1 OR 1 = 1)
					 AND ` + sqlTranNumFilter + `
				)
			) AS a
			ORDER BY created_date DESC
		`
	}

	// adminDeliverySearchFilters
	//adminOrderSearchFilters

	// DROPSHIPPER
	if userType == "Dropshipper" {
		sql = `
			SELECT 
			  IF(DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
			  d.name AS recipient,
			  d.id AS transaction_number,
			  0 AS amount,
			  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
			  do.name AS type,
			  ds.name AS status,
			  "" AS delivery_payment_method,
			  CONCAT(u.lastname, ', ', u.firstname) AS seller,
			  d.policy_number,
			  d.service_fee,
			  d.seller_id,
			  d.note,
			  d.address,
			  d.contact_number,
			  d.item_description
			FROM delivery d
			INNER JOIN delivery_status ds 
			  ON 1 = 1
				AND d.delivery_status_id = ds.id 
			INNER JOIN delivery_option do
			  ON 1 = 1 
				AND d.delivery_option_id = do.id
			INNER JOIN user u  
			  ON 1 = 1
				AND u.id = d.seller_id
			WHERE 1 = 1
			 AND d.dropshipper_id = ?
			 AND ` + sqlTranNumFilter + `
			 AND ` + sqlUserFilter + `
			ORDER BY d.created_date DESC
		`
	}

	// SELLER
	if userType == "Seller" {
		// Do nothing for now lol
		sqlUserFilter = "u.id = " + strconv.Itoa(userId)

		sql = `
			SELECT 
			  IF(DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
			  d.name AS recipient,
			  d.id AS transaction_number,
			  0 AS amount,
			  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
			  do.name AS type,
			  ds.name AS status,
			  "" AS delivery_payment_method,
			  CONCAT(u.lastname, ', ', u.firstname) AS seller,
			  d.policy_number,
			  d.service_fee,
			  d.seller_id,
			  d.note,
			  d.address,
			  d.contact_number,
			  d.item_description
			FROM delivery d
			INNER JOIN delivery_status ds 
			  ON 1 = 1
				AND d.delivery_status_id = ds.id 
			INNER JOIN delivery_option do
			  ON 1 = 1 
				AND d.delivery_option_id = do.id
			INNER JOIN user u  
			  ON 1 = 1
				AND u.id = d.seller_id
			WHERE 1 = 1
			 AND ` + sqlUserFilter + `
			 AND ` + sqlTranNumFilter + `
			ORDER BY d.created_date DESC
		`
	}

	// RIDER

	/**
	Fetches the transactions that are "Same Day Delivery"
	 */
	if userType == "Rider" {
		sql = `
			SELECT 
			  date_created,
			  recipient,
			  transaction_number,
			  CAST(amount AS DECIMAL(65,2)) AS amount,
			  tracking_number,
			  ` + "`type`" + `,  
			  ` + "`status`" + `,
			  address,
			  contact_number,
			  items,
			  delivery_payment_method
			FROM (
				(
					SELECT 
     				  d.created_date,
					  IF(DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
					  d.name AS recipient,
					  d.id AS transaction_number,
					  d.declared_amount AS amount,
					  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
					  do.name AS type,
					  ds.name AS status,
					  d.address,
					  d.contact_number,
 					  (SELECT COUNT(id) FROM delivery_detail WHERE delivery_id = d.id) AS items,
					  dpm.name AS delivery_payment_method
					FROM delivery d
					INNER JOIN delivery_status ds 
					  ON 1 = 1
						AND d.delivery_status_id = ds.id 
					INNER JOIN delivery_option do
					  ON 1 = 1 
						AND d.delivery_option_id = do.id
					INNER JOIN delivery_payment_method dpm
            		  ON 1 = 1
            		    AND d.delivery_payment_method_id = dpm.id
					WHERE 1 = 1
					 AND ds.name = ?
					 AND d.is_active = 1
					 AND d.region_id = ?
				)
			) AS a
			ORDER BY created_date DESC
		`
	}

	var err error
	var count int

	/**
	Pagination Logic
	 */

	paginate := func () (*[]ResponseTransactions, response_builder.Pagination, error) {
		user := UserModel.User{ID: userId}

		res, err := user.GetOne()

		if userType == "Admin" {
			count, err = database.GetQueryCount(sql)
		} else {
			if userType == "Rider" {
				count, err = database.GetQueryCount(sql, deliveryStatus, res[0].RegionId)
			} else if userType == "Seller" {
				count, err = database.GetQueryCount(sql)
			} else {
				count, err = database.GetQueryCount(sql, userId)
			}
		}

		// Fail error if error
		if err != nil {
			return &container, pagination, err
		}

		// Just return blank entries if there's no count (save operations)
		if count == 0 {
			return &container, pagination, nil
		}

		sql, pages, rowsPerPage, offset, page, totalCount, resultCount := database.GetPaginationDetails(
			sql,
			count,
			page,
			rows,
			1000,
		)

		pagination.SetData(rowsPerPage, offset, pages, rows, page, totalCount, resultCount)

		// Perform query
		if userType == "Admin" {
			fmt.Println("hkejhlkejhlekj")
			err = database.DBInstancePublic.Select(&container, sql)
		} else {
			fmt.Println("hkejhlkejhlekj")
			if userType == "Rider" {
				err = database.DBInstancePublic.Select(&container, sql, deliveryStatus, res[0].RegionId)
			} else if userType == "Seller" {
				fmt.Println("hkejhlkejhlekj")
				err = database.DBInstancePublic.Select(&container, sql)
			} else {
				err = database.DBInstancePublic.Select(&container, sql, userId)
			}
		}

		if err != nil {
			return &container, pagination, err
		}

		return &container, pagination, nil
	}

	res, pagination, err := paginate()

	return res, pagination, err
}

func (d *Delivery) GetServiceFee(orderDetails string) (*ResponseServiceFee, error) {
	var responseServiceFee ResponseServiceFee

	sql := `CALL get_service_fee(?)`

	err := database.DBInstancePublic.Get(&responseServiceFee, sql, orderDetails)

	return &responseServiceFee, err
}

func (d *Delivery) GetDeliveryDetails(userId int, userType string) (ResponseDeliveryDetails, error) {
	var responseDeliveryDetails ResponseDeliveryDetails

	var responseDeliveryDetailsItems []ResponseDeliveryDetailsItems
	var responseDeliveryDetailsInfo ResponseDeliveryDetailsInfo
	var responseDeliveryDetailsTracking []ResponseDeliveryDetailsTracking

	sqlDeliveryItems := ""
	sqlDeliveryInfo := ""
	sqlDeliveryTracking := ""

	if userType == "Seller" {
		sqlDeliveryItems = `
			SELECT
			  p.name AS product,
			  pt.name AS category,
			  dd.quantity
			FROM
			  delivery d
			  INNER JOIN delivery_detail dd
				ON 1 = 1
				AND d.id = dd.delivery_id
			  INNER JOIN product p 
				ON 1 = 1
				  AND dd.product_id = p.id
			  INNER JOIN product_type pt
				ON 1 = 1
				  AND p.product_type_id = pt.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.seller_id = ?
			  AND d.is_active = 1
			  AND p.is_active = 1
		`

		err := database.DBInstancePublic.Select(&responseDeliveryDetailsItems, sqlDeliveryItems, d.ID, userId)

		if err != nil {
			return responseDeliveryDetails, err
		}

		sqlDeliveryInfo = `
			SELECT
			  DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') AS date_created,
			  r.name AS region,
			  CONCAT(u.lastname, ', ', u.firstname) AS dropshipper,
			  d.address,
			  d.declared_amount,
			  d.service_fee,
			  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
			  d.contact_number,
			  d.note,
			  u2.mobile_number AS seller_mobile_number,
			  d.base_price,
			  d.name AS buyer_name,
			  u2.m88_account AS seller_m88_account,
			  CONCAT(u2.lastname, ', ', u2.firstname) AS seller_name,
			  IF(image_url IS NULL, '', image_url) AS image_url,
			  IF(item_description IS NULL, '', item_description) AS item_description
			FROM
			  delivery d
			  INNER JOIN region r 
				ON 1 = 1
				  AND d.region_id = r.id
			  INNER JOIN user u 
				ON 1 = 1
				  AND d.dropshipper_id = u.id
			  INNER JOIN user u2 
				ON 1 = 1
				  AND d.seller_id = u2.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.seller_id = ?
			  AND d.is_active = 1
			  AND u.is_active = 1	
		`

		err = database.DBInstancePublic.Get(&responseDeliveryDetailsInfo, sqlDeliveryInfo, d.ID, userId)

		fmt.Println(responseDeliveryDetailsInfo)

		if err != nil {
			return responseDeliveryDetails, err
		}

		sqlDeliveryTracking = `
			SELECT
			  DATE_FORMAT(
				CONVERT_TZ(dt.last_updated,'+00:00','+08:00'),
				'%Y-%m-%d %h:%i %p'
			  ) AS date_created,
			  ds.name AS status
			FROM
			  delivery d
			  INNER JOIN delivery_tracking dt
				ON 1 = 1
				AND d.id = dt.delivery_id
			  INNER JOIN delivery_status ds
				ON 1 = 1
				AND dt.delivery_status_id = ds.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.seller_id = ?
			  AND d.is_active = 1
			  AND d.is_active = 1		  
			ORDER BY d.created_date DESC
		`

		err = database.DBInstancePublic.Select(&responseDeliveryDetailsTracking, sqlDeliveryTracking, d.ID, userId)

		if err != nil {
			return responseDeliveryDetails, err
		}

		// Build output
		responseDeliveryDetails.Info = &responseDeliveryDetailsInfo
		responseDeliveryDetails.Items = &responseDeliveryDetailsItems
		responseDeliveryDetails.Tracking = &responseDeliveryDetailsTracking
	}

	if userType == "Dropshipper" {
		sqlDeliveryItems = `
			SELECT
			  p.name AS product,
			  pt.name AS category,
			  dd.quantity
			FROM
			  delivery d
			  INNER JOIN delivery_detail dd
				ON 1 = 1
				AND d.id = dd.delivery_id
			  INNER JOIN product p 
				ON 1 = 1
				  AND dd.product_id = p.id
			  INNER JOIN product_type pt
				ON 1 = 1
				  AND p.product_type_id = pt.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.dropshipper_id = ?
			  AND d.is_active = 1
			  AND p.is_active = 1
		`

		err := database.DBInstancePublic.Select(&responseDeliveryDetailsItems, sqlDeliveryItems, d.ID, userId)

		if err != nil {
			return responseDeliveryDetails, err
		}

		sqlDeliveryInfo = `
			SELECT
			  DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') AS date_created,
			  r.name AS region,
			  CONCAT(u.lastname, ', ', u.firstname) AS dropshipper,
			  d.address,
			  d.declared_amount,
			  d.service_fee,
			  d.contact_number,
			  d.note,
			  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
			  u2.mobile_number AS seller_mobile_number,
			  d.base_price,
			  d.name AS buyer_name,
			  u2.m88_account AS seller_m88_account,
			  CONCAT(u2.lastname, ', ', u2.firstname) AS seller_name,
			  IF(image_url IS NULL, '', image_url) AS image_url,
			  IF(item_description IS NULL, '', item_description) AS item_description
			FROM
			  delivery d
			  INNER JOIN region r 
				ON 1 = 1
				  AND d.region_id = r.id
			  INNER JOIN user u 
				ON 1 = 1
				  AND d.dropshipper_id = u.id
			  INNER JOIN user u2 
				ON 1 = 1
				  AND d.seller_id = u2.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.dropshipper_id = ?
			  AND d.is_active = 1
			  AND u.is_active = 1	
		`

		err = database.DBInstancePublic.Get(&responseDeliveryDetailsInfo, sqlDeliveryInfo, d.ID, userId)

		if err != nil {
			return responseDeliveryDetails, err
		}

		sqlDeliveryTracking = `
			SELECT
			  DATE_FORMAT(
				CONVERT_TZ(dt.last_updated,'+00:00','+08:00'),
				'%Y-%m-%d %h:%i %p'
			  ) AS date_created,
			  ds.name AS status
			FROM
			  delivery d
			  INNER JOIN delivery_tracking dt
				ON 1 = 1
				AND d.id = dt.delivery_id
			  INNER JOIN delivery_status ds
				ON 1 = 1
				AND dt.delivery_status_id = ds.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.dropshipper_id = ?
			  AND d.is_active = 1
			  AND d.is_active = 1		  
			ORDER BY d.created_date DESC
		`

		err = database.DBInstancePublic.Select(&responseDeliveryDetailsTracking, sqlDeliveryTracking, d.ID, userId)

		if err != nil {
			return responseDeliveryDetails, err
		}

		// Build output
		responseDeliveryDetails.Info = &responseDeliveryDetailsInfo
		responseDeliveryDetails.Items = &responseDeliveryDetailsItems
		responseDeliveryDetails.Tracking = &responseDeliveryDetailsTracking
	}

	if userType == "Admin" {
		sqlDeliveryItems = `
			SELECT
			  p.name AS product,
			  pt.name AS category,
			  dd.quantity
			FROM
			  delivery d
			  INNER JOIN delivery_detail dd
				ON 1 = 1
				AND d.id = dd.delivery_id
			  INNER JOIN product p 
				ON 1 = 1
				  AND dd.product_id = p.id
			  INNER JOIN product_type pt
				ON 1 = 1
				  AND p.product_type_id = pt.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.is_active = 1
			  AND p.is_active = 1
		`

		err := database.DBInstancePublic.Select(&responseDeliveryDetailsItems, sqlDeliveryItems, d.ID)

		if err != nil {
			return responseDeliveryDetails, err
		}

		sqlDeliveryInfo = `
			SELECT
			  DATE_FORMAT(CONVERT_TZ(d.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') AS date_created,
			  r.name AS region,
			  CONCAT(u.lastname, ', ', u.firstname) AS dropshipper,
			  d.address,
			  d.declared_amount,
			  d.service_fee,
			  d.contact_number,
			  d.note,
			  IF(d.tracking_number IS NULL, "", d.tracking_number) AS tracking_number,
			  u2.mobile_number AS seller_mobile_number,
			  d.base_price,
			  d.name AS buyer_name,
			  u2.m88_account AS seller_m88_account,
			  CONCAT(u2.lastname, ', ', u2.firstname) AS seller_name,
			  IF(image_url IS NULL, '', image_url) AS image_url,
			  IF(item_description IS NULL, '', item_description) AS item_description
			FROM
			  delivery d
			  INNER JOIN region r 
				ON 1 = 1
				  AND d.region_id = r.id
			  INNER JOIN user u 
				ON 1 = 1
				  AND d.dropshipper_id = u.id
			  INNER JOIN user u2 
				ON 1 = 1
				  AND d.seller_id = u2.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.is_active = 1
			  AND u.is_active = 1	
		`

		err = database.DBInstancePublic.Get(&responseDeliveryDetailsInfo, sqlDeliveryInfo, d.ID)

		if err != nil {
			return responseDeliveryDetails, err
		}

		sqlDeliveryTracking = `
			SELECT
			  DATE_FORMAT(
				CONVERT_TZ(dt.last_updated,'+00:00','+08:00'),
				'%Y-%m-%d %h:%i %p'
			  ) AS date_created,
			  ds.name AS status
			FROM
			  delivery d
			  INNER JOIN delivery_tracking dt
				ON 1 = 1
				AND d.id = dt.delivery_id
			  INNER JOIN delivery_status ds
				ON 1 = 1
				AND dt.delivery_status_id = ds.id
			WHERE 1 = 1
			  AND d.id = ?
			  AND d.is_active = 1
			  AND d.is_active = 1		  
			ORDER BY d.created_date DESC
		`

		err = database.DBInstancePublic.Select(&responseDeliveryDetailsTracking, sqlDeliveryTracking, d.ID)

		if err != nil {
			return responseDeliveryDetails, err
		}

		// Build output
		responseDeliveryDetails.Info = &responseDeliveryDetailsInfo
		responseDeliveryDetails.Items = &responseDeliveryDetailsItems
		responseDeliveryDetails.Tracking = &responseDeliveryDetailsTracking
	}

	if sqlDeliveryInfo == "" || sqlDeliveryItems == "" || sqlDeliveryTracking == "" {
		return responseDeliveryDetails, errors.New("no sql ran for GetDeliveryDetails")
	}

	return responseDeliveryDetails, nil
}

func Top10Sellers() (*[]Top10Seller, error) {
	var Top10Sellers []Top10Seller

	sql := `
		SELECT
		  d.seller_id,
		  CONCAT(u.lastname, ', ', u.firstname) AS seller,
		  COUNT(*) AS deliveries,
		  SUM(d.declared_amount) AS sales
		FROM
		  delivery d
		  INNER JOIN user u
			ON 1 = 1
			AND d.seller_id = u.id
		  INNER JOIN delivery_status ds
			 ON 1 = 1
			   AND d.delivery_status_id = ds.id
			   AND ds.name IN ('Fulfilled', 'Delivered')
		WHERE 1 = 1
		  AND d.created_date BETWEEN '2020-09-01' AND  '2020-09-30'
		GROUP BY u.id
		ORDER BY sales DESC, deliveries DESC
		LIMIT 10
	`

	err := database.DBInstancePublic.Select(&Top10Sellers, sql)

	return &Top10Sellers, err
}

func (d *Delivery) GetDashboardDeliveryStatus(userId int, userType string) (*ResponseDashboardDeliveryStatus, error) {
	var responseDashboardDeliveryStatusContainer []ResponseDashboardDeliveryStatusContainer
	var responseDashboardDeliveryStatus ResponseDashboardDeliveryStatus

	var sql string

	if userType == "Seller" {
		sql = `
			SELECT
			  ds.name,
			  IF (a.amount IS NULL, 0, a.amount) AS amount
			FROM
			  delivery_status ds
			  LEFT JOIN
				(SELECT
				  IF(COUNT(*) IS NULL, 0, COUNT(*)) AS amount,
				 
				  ds.name AS delivery_status
				FROM
				  delivery d
				  INNER JOIN delivery_status ds
					ON 1 = 1
					AND d.delivery_status_id = ds.id
				WHERE 1 = 1
				  AND d.is_active = 1
				  AND d.seller_id = ?
				  GROUP BY ds.name) a
				ON 1 = 1
				AND a.delivery_status = ds.name
		`
	}

	if userType == "Dropshipper" {
		sql = `
			SELECT
			  ds.name,
			  IF (a.amount IS NULL, 0, a.amount) AS amount
			FROM
			  delivery_status ds
			  LEFT JOIN
				(SELECT
				  IF(COUNT(*) IS NULL, 0, COUNT(*)) AS amount,
				  ds.name AS delivery_status
				FROM
				  delivery d
				  INNER JOIN delivery_status ds
					ON 1 = 1
					AND d.delivery_status_id = ds.id
				WHERE 1 = 1
				  AND d.is_active = 1
				  AND d.dropshipper_id = ?
				  GROUP BY ds.name) a
				ON 1 = 1
				AND a.delivery_status = ds.name
		`
	}

	if sql == "" {
		return &responseDashboardDeliveryStatus, errors.New("no valid user type to run GetDashboardDeliveryStatus")
	}

	// Run lol
	err := database.DBInstancePublic.Select(
		&responseDashboardDeliveryStatusContainer,
		sql,
		userId,
	)

	if err != nil {
		return &responseDashboardDeliveryStatus, err
	}

	/**
		Transform values
	 */
	for _, v := range responseDashboardDeliveryStatusContainer {
		// responseDashboardDeliveryStatus
		if v.Name == "Proposed" {
			responseDashboardDeliveryStatus.Proposed = v.Amount
		}

		if v.Name == "Accepted" {
			responseDashboardDeliveryStatus.Accepted = v.Amount
		}

		if v.Name == "Fulfilled" {
			responseDashboardDeliveryStatus.Fulfilled = v.Amount
		}

		if v.Name == "Rejected" {
			responseDashboardDeliveryStatus.Rejected = v.Amount
		}

		if v.Name == "Delivered" {
			responseDashboardDeliveryStatus.Delivered = v.Amount
		}
	}

	var (
		DeclaredAmount float64
	)

	// Get total sales
	if userType == "Seller" {
		sql = `
			SELECT 
				IF(SUM(declared_amount) IS NULL, 0, SUM(declared_amount))
			FROM delivery d 
			INNER JOIN delivery_status ds 
				ON 1 = 1
					AND d.delivery_status_id = ds.id 
					AND ds.name = 'Delivered'
			WHERE  1 = 1
				AND seller_id = ? 
				AND is_active = 1
		`

		err := database.DBInstancePublic.Get(&DeclaredAmount, sql, userId)

		if err != nil {
			return &responseDashboardDeliveryStatus, err
		}

		responseDashboardDeliveryStatus.TotalSales = DeclaredAmount
	}

	if userType == "Dropshipper" {
		// Old
		/*sql = `
			SELECT 
				IF(SUM(declared_amount) IS NULL, 0, SUM(declared_amount))
			FROM delivery d 
			INNER JOIN delivery_status ds 
				ON 1 = 1
					AND d.delivery_status_id = ds.id 
					AND ds.name IN ('Fulfilled', 'Delivered')
			WHERE  1 = 1
				AND d.dropshipper_id = ? 
				AND is_active = 1
		`*/

		sql = `
			SELECT 
				0 AS declared_amount
			FROM delivery d 
			INNER JOIN delivery_status ds 
				ON 1 = 1
					AND d.delivery_status_id = ds.id 
					AND ds.name IN ('Fulfilled', 'Delivered')
			WHERE  1 = 1
				AND d.dropshipper_id = ? 
				AND is_active = 1
		`

		err := database.DBInstancePublic.Get(&DeclaredAmount, sql, userId)

		if err != nil {
			return &responseDashboardDeliveryStatus, err
		}

		responseDashboardDeliveryStatus.TotalSales = DeclaredAmount
	}


	return &responseDashboardDeliveryStatus, nil

	// Perform different queries whether this is a seller or a dropshipper

}