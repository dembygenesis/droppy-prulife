package coin_transactions

import (
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func (t CoinTransaction) GetAll() (*[]ResponseCoinTransactionList, error) {
	var responseTransactionList []ResponseCoinTransactionList

	sql := `
		SELECT
		  ct.id,
          IF(DATE_FORMAT(CONVERT_TZ(ct.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p') IS NULL, "", DATE_FORMAT(CONVERT_TZ(ct.created_date,'+00:00','+08:00'), '%Y-%m-%d %h:%i %p')) AS date_created,
          IF(CONCAT(u2.lastname, ', ', u2.firstname) IS NULL, "", CONCAT(u2.lastname, ', ', u2.firstname)) AS created_by,
          IF(CONCAT(u.lastname, ', ', u.firstname) IS NULL, "", CONCAT(u.lastname, ', ', u.firstname)) AS allotted_to,
          IF(CONCAT(u3.lastname, ', ', u3.firstname) IS NULL, "", CONCAT(u3.lastname, ', ', u3.firstname)) AS updated_by,
          ct.is_active,
          ct.type,
          ct.amount
		FROM
		  coin_transaction ct
		  INNER JOIN user u
			ON 1 = 1
			AND ct.user_id = u.id
		  LEFT JOIN user u2
			ON 1 = 1
			AND ct.created_by = u2.id
		  LEFT JOIN user u3
			ON 1 = 1
			AND ct.updated_by = u3.id
		ORDER BY ct.created_date DESC,
		  ct.type ASC
	`

	err := database.DBInstancePublic.Select(&responseTransactionList, sql)

	return &responseTransactionList, err
}

/*func (t CoinTransaction) Create() (sql.Result, error) {
	var sqlResult sql.Result

	sql := `
		SELECT 5 as hahah
	`

	fmt.Println("sqlResult", sqlResult)

	err := database.DBInstancePublic.Select(&sqlResult, sql)

	return sqlResult, err
}*/