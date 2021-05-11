package delivery

type DeliveryOption struct {
	Id   uint
	Name string
}

type ParamsDeliveryCreate struct {
	// Insert fields
	Id               uint
	CreatedBy        int
	UpdatedBy        int
	CreatedDate      string
	IsActive         int
	Name             string
	Address          string
	ServiceFee       float64
	DeliveryOptionId int
	DeliveryStatusId int
	SellerId         int
	DropshipperId    int
	ContactNumber    string
	Note             string
	ImageUrl         string
	ItemDescription  string
	PolicyNumber     string

	// Display fields
}
