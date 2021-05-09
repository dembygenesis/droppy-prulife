package coin_transactions

// Transaction is an exact map of the schema
type CoinTransaction struct {
	ID int `json:"id" db:"id"`
}

// ResponseTransactionList is an exact map of the transaction list query
type ResponseCoinTransactionList struct {
	ID          int     `json:"id" db:"id"`
	DateCreated string  `json:"date_created" db:"date_created"`
	CreatedBy   string  `json:"created_by" db:"created_by"`
	AllottedTo  string  `json:"allotted_to" db:"allotted_to"`
	UpdatedBy   string  `json:"updated_by" db:"updated_by"`
	IsActive    int     `json:"is_active" db:"is_active"`
	Type        string  `json:"type" db:"type"`
	Amount      float64 `json:"amount" db:"amount"`
}
