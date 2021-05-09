package delivery

type CreateDelivery struct {
	DeliveryOption        string  `json:"delivery_option" form:"delivery_option" db:"delivery_option"`
	SellerId              int     `json:"seller_id" form:"seller_id" db:"seller_id"`
	DropshipperId         int     `json:"dropshipper_id" form:"dropshipper_id" db:"dropshipper_id"`
	Name                  string  `json:"name" form:"name" db:"name"`
	ContactNumber         string  `json:"contact_number" form:"contact_number" db:"contact_number"`
	Address               string  `json:"address" form:"address" db:"address"`
	Note                  string  `json:"note" form:"note" db:"note"`
	Region                string  `json:"region" form:"region" db:"region"`
	ServiceFee            float64 `json:"service_fee" form:"service_fee" db:"service_fee" required:"false"`
	DeclaredAmount        float64 `json:"declared_amount" form:"declared_amount" db:"declared_amount"`
	DeliveryDetails       string  `json:"delivery_details" form:"delivery_details" db:"delivery_details" required:"false"`

	// We need new fields: Picture, and Item Description
	ImageUrl        string `json:"image_url" form:"image_url" db:"image_url" required:"false"`
	ItemDescription string `json:"item_description" form:"item_description" db:"item_description"`
}

// I can even do the validations here... But I'll leave that all in my repository.