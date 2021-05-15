package withdrawals

type Withdrawal struct {
	ID       int    `json:"id"`
	UserType string `json:"user_type"`
	UserID   int    `json:"user_id"`
}

/**
Response
*/
type ResponseWithdrawalList struct {
	ID                 float64 `json:"id,omitempty" db:"id"`
	DateCreated        string  `json:"date_created,omitempty" db:"date_created"`
	Amount             float64 `json:"amount,omitempty" db:"amount"`
	BankNo             string  `json:"bank_no,omitempty" db:"bank_no"`
	BankName           string  `json:"bank_name,omitempty" db:"bank_name"`
	BankAccountName    string  `json:"bank_account_name,omitempty" db:"bank_account_name"`
	Fee                float64 `json:"fee,omitempty" db:"fee"`
	TotalAmount        float64 `json:"total_amount,omitempty" db:"total_amount"`
	WithdrawalStatus   string  `json:"withdrawal_status,omitempty" db:"withdrawal_status"`
	VoidOrRejectReason string  `json:"void_or_reject_reason,omitempty" db:"void_or_reject_reason"`
	UserFullName       string  `json:"user_fullname,omitempty" db:"user_fullname"`
	UserType           string  `json:"user_type,omitempty" db:"user_type"`
	Description        string  `json:"description,omitempty" db:"description"`
	ContactNo          string  `json:"contact_no,omitempty" db:"contact_no"`
}

/**
Params
*/
type ParamsCreateWithdrawal struct {
	Amount          float64 `json:"amount" db:"amount"`
	BankNo          string  `json:"bank_no" db:"bank_no"`
	BankTypeId      float64 `json:"bank_type_id" db:"bank_type_id"`
	BankAccountName string  `json:"bank_account_name" db:"bank_account_name"`
	ContactNo       string  `json:"contact_no" db:"contact_no"`
}

type ParamsUpdateWithdrawal struct {
	ID                 float64 `json:"id" db:"id"`
	Status             string  `json:"status" db:"status"`
	VoidOrRejectReason string  `json:"void_or_reject_reason" db:"void_or_reject_reason"`
	ReferenceNumber    string  `json:"reference_number" db:"reference_number"`
}
