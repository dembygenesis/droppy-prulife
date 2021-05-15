package withdrawals

import (
	"database/sql"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	string2 "github.com/dembygenesis/droppy-prulife/utilities/string"
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
		CALL add_withdrawal(?, ?, ?, ?, ?, ?)
	`

	res, err := database.DBInstancePublic.Exec(sql,
		w.UserID,
		p.Amount,
		p.BankNo,
		p.BankTypeId,
		p.BankAccountName,
		p.ContactNo,
	)

	return res, err
}

func (w *Withdrawal) GetAll(page int, rows int, filters []string) (*[]ResponseWithdrawalList, response_builder.Pagination, error) {
	var container []ResponseWithdrawalList
	var pagination response_builder.Pagination

	// Handle filters
	filter := ""

	isDropshipperOrSeller := w.UserType == "Dropshipper" || w.UserType == "Seller"

	if isDropshipperOrSeller {
		filter += "AND w.user_id = ? "
	}

	if len(filters) != 0 {
		filter += "AND ws.name IN (" + string2.GetJoinedStringForWhereIn(filters) + ") "
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
		  END AS description,
		  IF(w.contact_no IS NULL, "", w.contact_no) AS contact_no
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
		  ` + filter + `
		ORDER BY w.created_date DESC
	`

	/**
	Pagination Logic
	*/

	paginate := func () (*[]ResponseWithdrawalList, response_builder.Pagination, error) {
		var count int
		var err error

		/**
		Use this
		*/
		if isDropshipperOrSeller {
			fmt.Println("Executing here becsause dropshipper or sellers", w.UserID)
			count, err = database.GetQueryCount(sql, w.UserID)
		} else {
			count, err = database.GetQueryCount(sql)
		}

		fmt.Println("================= withdrawal", sql)
		// End this

		// Fail error if error
		if err != nil {
			return &container, pagination, err
		}

		// Just return blank entries if there's no count (save operations)
		if count == 0 {
			pagination.Pages = make([]int, 0)
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
		if isDropshipperOrSeller {
			err = database.DBInstancePublic.Select(&container, sql, w.UserID)
		} else {
			err = database.DBInstancePublic.Select(&container, sql)
		}

		if err != nil {
			return &container, pagination, err
		}

		return &container, pagination, nil
	}

	res, pagination, err := paginate()

	return res, pagination, err
}