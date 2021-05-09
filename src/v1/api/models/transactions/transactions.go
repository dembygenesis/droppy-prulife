package transactions

import (
	"database/sql"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (t *Transaction) GetAll() (*[]ResponseTransactionList, error) {
	var responseTransactionList []ResponseTransactionList

	sql := `
		SELECT
		  t.id,
		  IF(DATE_FORMAT(CONVERT_TZ(t.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(t.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
		  IF(CONCAT(u.lastname, ', ', u.firstname) IS NULL, '', CONCAT(u.lastname, ', ', u.firstname)) AS created_by,
		  IF(CONCAT(u2.lastname, ', ', u2.firstname) IS NULL, '', CONCAT(u2.lastname, ', ', u2.firstname)) AS updated_by,
		  IF(CONCAT(u3.lastname, ', ', u3.firstname) IS NULL, '', CONCAT(u3.lastname, ', ', u3.firstname)) AS admin_allotted,
		  IF(CONCAT(u4.lastname, ', ', u4.firstname) IS NULL, '', CONCAT(u4.lastname, ', ', u4.firstname)) AS user_allotted,
		  IF(IF(bt.name IS NULL, '', bt.name) IS NULL, '', IF(bt.name IS NULL, '', bt.name)) AS bank,
		  t.is_active,
		  t.reference_number,
		  t.description,
		  t.amount,
		  t.coin_amount,
		  t.money_in
		FROM
		  transaction t
		  INNER JOIN` + "`user`" + `u 
			ON 1 = 1
			  AND t.created_by = u.id
		  LEFT JOIN` + "`user`" + `u2
			ON 1 = 1
			  AND t.updated_by = u2.id
		  INNER JOIN` + "`user`" + `u3 
			ON 1 = 1
			  AND t.admin_allotted_id = u3.id
		  INNER JOIN` + "`user`" + `u4
			ON 1 = 1
			  AND t.user_allotted_id = u4.id
		  INNER JOIN bank_type bt 
			ON 1 = 1
			  AND t.bank_type_id = bt.id
		  
		  ORDER BY t.created_date DESC
	`

	fmt.Println(sql)

	err := database.DBInstancePublic.Select(&responseTransactionList, sql)

	return &responseTransactionList, err
}

func (t *Transaction) Create(p *ParamsTransaction) (*sql.Result, error) {
	sql := `
		CALL add_transaction(?, ?, ?, ?, ?, ?, ?, ?)
	`

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		p.Amount,
		p.CoinAmount,
		p.AdminId,
		p.UserId,
		p.MoneyIn,
		p.BankTypeId,
		p.ReferenceNumber,
		p.Description,
	)

	return &sqlResult, err
}

func (t *Transaction) Delete(p *ParamsTransactionDelete) (*sql.Result, error) {
	sql := `
		CALL void_transaction(?, ?)
	`

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		p.ID,
		p.AdminId,
	)

	return &sqlResult, err
}