package products

type Product struct {
	ID            int    `json:"id,omitempty" db:"id"`
	Name          string `json:"name,omitempty" db:"name"`
	ProductTypeId int    `json:"product_type_id,omitempty" db:"product_type_id"`
	CreatedBy     int    `json:"created_by,omitempty" db:"created_by"`
	CreatedDate   string `json:"created_date,omitempty" db:"created_date"`
	LastUpdated   string `json:"last_updated,omitempty" db:"last_updated"`
	UpdatedBy     int    `json:"updated_by,omitempty" db:"updated_by"`
	IsActive      int    `json:"is_active,omitempty" db:"is_active"`
	Url           string `json:"url,omitempty" db:"url"`
}

/**
Params
*/

type ParamsUpdate struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	ProductTypeId int    `json:"productTypeId" db:"product_type_id"`
}

type ParamsDelete struct {
	ID int `json:"id" db:"id"`
}

type ParamsCreate struct {
	Name          string `json:"name" db:"name"`
	ProductTypeId int    `json:"product_type_id" db:"product_type_id"`
}

/**
Responses
*/
type ResponseProductList struct {
	ID                      int     `json:"id,omitempty" db:"id"`
	Name                    string  `json:"name,omitempty" db:"name"`
	Url                     string  `json:"url,omitempty" db:"url"`
	Category                string  `json:"category,omitempty" db:"category"`
	CategoryID              int     `json:"category_id,omitempty" db:"category_id"`
	PricePerItem            float64 `json:"price_per_item,omitempty" db:"price_per_item"`
	PricePerItemDropshipper float64 `json:"price_per_item_dropshipper,omitempty" db:"price_per_item_dropshipper"`
}

type ResponseProductSpecific struct {
	ID            int    `json:"id,omitempty" db:"id"`
	Name          string `json:"name,omitempty" db:"name"`
	Url           string `json:"url,omitempty" db:"url"`
	Category      string `json:"category,omitempty" db:"category"`
	ProductTypeId int    `json:"product_type_id,omitempty" db:"product_type_id"`
}

type ResponseInventoryList struct {
	ID        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Category  string  `json:"category" db:"category"`
	Remaining float64 `json:"remaining" db:"remaining"`
	Region    string  `json:"region" db:"region"`
}

type ResponseSellerList struct {
	SellerId   int    `json:"id" db:"id"`
	SellerName string `json:"name" db:"name"`
}
