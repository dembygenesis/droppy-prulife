package orders

type Order struct {
}

type ParamsOrder struct {
	// Admin id
	AdminId      int    `json:"admin_id"`
	UserId       int    `json:"user_id"`
	OrderDetails string `json:"order_details"`
}

type ParamsOrderUpdate struct {
	AdminId     int
	OrderId     int    `json:"order_id"`
	Description string `json:"void_or_reject_reason"`
}

type ParamsGetOrderDetails struct {
	OrderId int `json:"order_id" db:"order_id"`
}

type ResponseOrders struct {
	ID              int     `json:"id" db:"id"`
	Amount          float64 `json:"amount" db:"amount"`
	DateCreated     string  `json:"date_created" db:"date_created"`
	Seller          string  `json:"seller" db:"seller"`
	Product         string  `json:"product" db:"product"`
	PricePerItem    float64 `json:"price_per_item" db:"price_per_item"`
	Quantity        float64 `json:"quantity" db:"quantity"`
	TotalPrice      float64 `json:"total_price" db:"total_price"`
	IsActive        int     `json:"is_active" db:"is_active"`
	Region          string  `json:"region" db:"region"`
	OrderTotalPrice float64 `json:"order_total_price" db:"order_total_price"`
	Dropshipper     string  `json:"dropshipper" db:"dropshipper"`
}

type ResponseOrdersDisplay struct {
	ID              int     `json:"id" db:"id"`
	Amount          float64 `json:"amount" db:"amount"`
	DateCreated     string  `json:"date_created" db:"date_created"`
	Seller          string  `json:"seller" db:"seller"`
	Product         string  `json:"product" db:"product"`
	PricePerItem    float64 `json:"price_per_item" db:"price_per_item"`
	Quantity        float64 `json:"quantity" db:"quantity"`
	TotalPrice      float64 `json:"total_price" db:"total_price"`
	IsActive        int     `json:"is_active" db:"is_active"`
	Region          string  `json:"region" db:"region"`
	OrderTotalPrice float64 `json:"order_total_price" db:"order_total_price"`
	Dropshipper     string  `json:"dropshipper" db:"dropshipper"`
}

type ResponseDeliveryDetails struct {
	ProductDetails []struct {
	}
}
