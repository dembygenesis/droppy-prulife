package transactions

// Transaction - the fields here map to the transaction schema
type Transaction struct {
	ID string `json:"id"`
}

/**
Response Types
*/

// Transaction - the fields here are the generic all responses
type ResponseTransactionList struct {
	ID              int     `json:"id" db:"id"`
	DateCreated     string  `json:"date_created" db:"date_created"`
	CreatedBy       string  `json:"created_by" db:"created_by"`
	UpdatedBy       string  `json:"updated_by" db:"updated_by"`
	AdminAllotted   string  `json:"admin_allotted" db:"admin_allotted"`
	UserAllotted    string  `json:"user_allotted" db:"user_allotted"`
	Bank            string  `json:"bank" db:"bank"`
	IsActive        int     `json:"is_active" db:"is_active"`
	ReferenceNumber string  `json:"reference_number" db:"reference_number"`
	MoneyIn         int     `json:"money_in" db:"money_in"`
	Description     string  `json:"description" db:"description"`
	Amount          float64 `json:"amount" db:"amount"`
	CoinAmount      float64 `json:"coin_amount" db:"coin_amount"`
}

/**
Param Types

p_amount DECIMAL(15,2),
    p_coin_amount DECIMAL(15,2),
    p_admin_allotted_id INTEGER,
    p_user_allotted_id INTEGER,
    p_money_in BOOLEAN,
    p_bank_type_id INTEGER,
    p_reference_number TEXT,
    p_description TEXT

p_transaction_id
p_admin_id

*/
type ParamsTransaction struct {
	Amount          int    `json:"amount" db:"amount"`
	CoinAmount      int    `json:"coin_amount" db:"coin_amount"`
	AdminId         int    `json:"admin_id" db:"admin_id"`
	UserId          int    `json:"user_id" db:"user_id"`
	MoneyIn         bool   `json:"money_in" db:"money_in"`
	BankTypeId      int    `json:"bank_type_id" db:"bank_type_id"`
	ReferenceNumber string `json:"reference_number" db:"reference_number"`
	Description     string `json:"description" db:"description"`
}

type ParamsTransactionDelete struct {
	ID      int `json:"id" db:"id"`
	AdminId int `json:"admin_id" db:"admin_id"`
}
