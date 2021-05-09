package withdrawals

import (
	"database/sql"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (w *Withdrawal) Update(p ParamsUpdateWithdrawal) (sql.Result, error) {
	sql := `
		CALL update_withdrawal(?, ?, ?, ?, ?)
	`

	res, err := database.DBInstancePublic.Exec(sql,
		w.UserID,
		p.ID,
		p.Status,
		p.ReferenceNumber,
		p.VoidOrRejectReason,
	)

	return res, err
}

func (w *Withdrawal) Create(p ParamsCreateWithdrawal) (sql.Result, error) {
	sql := `
		CALL add_withdrawal(?, ?, ?, ?, ?)
	`

	res, err := database.DBInstancePublic.Exec(sql,
		w.UserID,
		p.Amount,
		p.BankNo,
		p.BankTypeId,
		p.BankAccountName,
	)

	return res, err
}

func (w *Withdrawal) GetAll() (*[]ResponseWithdrawalList, error) {
	var responseWithdrawalList []ResponseWithdrawalList

	// Default is all
	filter := "1 = 1"

	isDropshipperOrSeller := w.UserType == "Dropshipper" || w.UserType == "Seller"

	if isDropshipperOrSeller {
		filter = "w.user_id = ?"
	}

	sql := `
		SELECT
          w.id,
		  DATE_FORMAT(CONVERT_TZ(w.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') AS date_created,
		  w.amount,
		  w.bank_no,
		  bt.name AS bank_name,
		  w.bank_account_name,
		  w.fee,
		  w.total_amount,
		  ws.name AS withdrawal_status,
		  IF(w.void_or_reject_reason IS NULL, "", w.void_or_reject_reason) AS void_or_reject_reason,
		  CONCAT(u.lastname, ', ', u.firstname) AS user_fullname,
		  ut.name AS user_type,
		  CASE 
			WHEN ws.name = 'Pending' THEN 'N/A'
			WHEN ws.name = 'Completed' THEN IF(w.reference_number IS NULL, '', w.reference_number)
			WHEN ws.name = 'Voided' THEN IF(w.void_or_reject_reason IS NULL, '', w.void_or_reject_reason) 
		  END AS description
		FROM
		  withdrawal w
		  INNER JOIN withdrawal_status ws
			ON 1 = 1
			AND w.withdrawal_status_id = ws.id
		  INNER JOIN ` + "`user`" + ` u
			ON 1 = 1
			AND w.user_id = u.id
 		  INNER JOIN user_type ut
			ON 1 = 1
			AND u.user_type_id = ut.id
	 	  INNER JOIN bank_type bt
			ON 1 = 1
			AND w.bank_type_id = bt.id
		WHERE 1 = 1
		  AND w.is_active = 1
		  AND ` + filter + `
		ORDER BY w.created_date DESC
	`

	var err error

	if isDropshipperOrSeller {
		err = database.DBInstancePublic.Select(&responseWithdrawalList, sql, w.UserID)
	} else {
		err = database.DBInstancePublic.Select(&responseWithdrawalList, sql)
	}

	return &responseWithdrawalList, err
}