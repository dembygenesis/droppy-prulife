package delivery_option

type DeliveryOption struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
