package delivery

type RequestCreateDelivery struct {
	CreatedById       int
	CreatedByUserType string
	DeliveryOption    string  `json:"delivery_option" form:"delivery_option" db:"delivery_option"`
	SellerId          int     `json:"seller_id" form:"seller_id" db:"seller_id"`
	DropshipperId     int     `json:"dropshipper_id" form:"dropshipper_id" db:"dropshipper_id"`
	Name              string  `json:"name" form:"name" db:"name"`
	ContactNumber     string  `json:"contact_number" form:"contact_number" db:"contact_number"`
	Address           string  `json:"address" form:"address" db:"address"`
	Note              string  `json:"note" form:"note" db:"note"`
	Region            string  `json:"region" form:"region" db:"region" required:"false"`
	ServiceFee        float64 `json:"service_fee" form:"service_fee" db:"service_fee" required:"false"`
	DeclaredAmount    float64 `json:"declared_amount" form:"declared_amount" db:"declared_amount" required:"false"`
	DeliveryDetails   string  `json:"delivery_details" form:"delivery_details" db:"delivery_details" required:"false"`
	ServiceFeeType    string  `json:"service_fee_type" form:"service_fee_type" db:"service_fee_type" required:"true"`
	PolicyNumber      string  `json:"policy_number" form:"policy_number" db:"policy_number" required:"true"`

	// We need new fields: Picture, and Item Description
	ImageUrl        string `json:"image_url" form:"image_url" db:"image_url" required:"false"`
	ItemDescription string `json:"item_description" form:"item_description" db:"item_description"`
}

type RequestUpdateDelivery struct {
	DeliveryId         int    `json:"delivery_id" form:"delivery_id" db:"delivery_id" required:"true"`
	DeliveryStatus     string `json:"delivery_status" form:"delivery_status" db:"delivery_status" required:"true"`
	VoidOrRejectReason string `json:"void_or_reject_reason" form:"void_or_reject_reason" db:"void_or_reject_reason" required:"false"`
	CreatedByUserType  string

	// Edit fields
	PolicyNumber    string `json:"policy_number" form:"policy_number" db:"policy_number" required:"false"`
	Name            string `json:"name" form:"name" db:"name" required:"false"`
	ContactNo       string `json:"contact_no" form:"contact_no" db:"contact_no" required:"false"`
	Note            string `json:"note" form:"note" db:"note" required:"false"`
	Address         string `json:"address" form:"address" db:"address" required:"false"`
	ItemDescription string `json:"item_description" form:"item_description" db:"item_description" required:"false"`
}
