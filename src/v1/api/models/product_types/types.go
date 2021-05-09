package product_types

type ProductType struct {
	ID   int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
