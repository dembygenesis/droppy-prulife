package deliveries

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Delivery struct {
	ID     int
	UserId int
}

type ResponseDashboardDeliveryStatusContainer struct {
	Name   string  `json:"name" db:"name"`
	Amount float64 `json:"amount" db:"amount"`
}

type ResponseDashboardDeliveryStatus struct {
	TotalSales float64 `json:"total_sales"`
	Accepted   float64 `json:"accepted"`
	Delivered  float64 `json:"delivered"`
	Fulfilled  float64 `json:"fulfilled"`
	Proposed   float64 `json:"proposed"`
	Rejected   float64 `json:"rejected"`
}

type ResponseTransactions struct {
	SellerId              int     `json:"seller_id" db:"seller_id"`
	DateCreated           string  `json:"date_created" db:"date_created"`
	Recipient             string  `json:"recipient" db:"recipient"`
	TransactionNumber     string  `json:"transaction_number" db:"transaction_number"`
	Amount                float64 `json:"amount" db:"amount"`
	TrackingNumber        string  `json:"tracking_number" db:"tracking_number"`
	Type                  string  `json:"type" db:"type"`
	Status                string  `json:"status" db:"status"`
	Items                 float64 `json:"items" db:"items"`
	DeliveryPaymentMethod string  `json:"delivery_payment_method" db:"delivery_payment_method"`

	Dropshipper     string `json:"dropshipper" db:"dropshipper"`
	SellerName      string `json:"seller" db:"seller"`
	Region          string `json:"region" db:"region"`
	Address         string `json:"address" db:"address"`
	ContactNumber   string `json:"contact_number" db:"contact_number"`
	ItemDescription string `json:"item_description" db:"item_description"`

	/*Add these new fields*/
	PolicyNumber string `json:"policy_number" db:"policy_number"`
	ServiceFee   string `json:"service_fee" db:"service_fee"`
	Note         string `json:"note" db:"note"`
}

/**
Response
*/

type ResponseDeliveryDetails struct {
	Items    *[]ResponseDeliveryDetailsItems    `json:"items"`
	Info     *ResponseDeliveryDetailsInfo       `json:"details"`
	Tracking *[]ResponseDeliveryDetailsTracking `json:"tracking"`
}

type ResponseDeliveryDetailsTracking struct {
	DateCreated string `json:"date_created" db:"date_created"`
	Status      string `json:"status" db:"status"`
}

type ResponseDeliveryDetailsItems struct {
	Product  string `json:"product" db:"product"`
	Category string `json:"category" db:"category"`
	Quantity int    `json:"quantity" db:"quantity"`
}

type ResponseDeliveryDetailsInfo struct {
	DateCreated        string  `json:"date_created" db:"date_created" `
	Region             string  `json:"region" db:"region"`
	Dropshipper        string  `json:"dropshipper" db:"dropshipper"`
	Address            string  `json:"address" db:"address"`
	DeclaredAmount     float64 `json:"declared_amount" db:"declared_amount"`
	ServiceFee         float64 `json:"service_fee" db:"service_fee"`
	TrackingNumber     string  `json:"tracking_number" db:"tracking_number"`
	ContactNumber      string  `json:"contact_number" db:"contact_number"`
	Note               string  `json:"note" db:"note"`
	SellerMobileNumber string  `json:"seller_mobile_number" db:"seller_mobile_number"`
	BasePrice          float64 `json:"base_price" db:"base_price"`
	SellerM88Account   string  `json:"seller_m88_account" db:"seller_m88_account"`
	SellerName         string  `json:"seller_name" db:"seller_name"`
	BuyerName          string  `json:"buyer_name" db:"buyer_name"`
	ImageUrl           string  `json:"image_url" db:"image_url"`
	ItemDescription    string  `json:"item_description" db:"item_description"`
}

type ResponseCoinTransactions struct {
	DateCreated        string  `json:"date_created" db:"date_created" `
	Type               string  `json:"type" db:"type"`
	Amount             float64 `json:"amount" db:"amount"`
	ReferenceNumber    string  `json:"reference_number" db:"reference_number"`
	BankType           string  `json:"bank_type" db:"bank_type"`
	Source             string  `json:"source" db:"source"`
	Recipient          string  `json:"recipient" db:"recipient"`
	TranNum            float64 `json:"tran_num" db:"tran_num"`
	VoidOrRejectReason string  `json:"void_or_reject_reason" db:"void_or_reject_reason"`
	SourceType         string  `json:"source_type" db:"source_type"`
	IsActive           float64 `json:"is_active" db:"is_active"`
	RunningBalance     float64 `json:"running_balance" db:"running_balance"`
	CreatedDate        string  `json:"created_date" db:"created_date"`
}

type ResponseServiceFee struct {
	ServiceFee float64 `json:"service_fee" db:"service_fee" `
}

type ResponseTestTransaction struct {
	Balance string `json:"firstname" db:"firstname"`
}

/**
Params
*/

type ParamsUpdateDelivery struct {
	DeliveryId         int    `json:"delivery_id" db:"delivery_id"`
	DeliveryStatus     string `json:"delivery_status" db:"delivery_status"`
	TrackingNumber     string `json:"tracking_number" db:"tracking_number"`
	VoidOrRejectReason string `json:"void_or_reject_reason" db:"void_or_reject_reason"`
}

type ParamsCreateOrder struct {
	SellerId      int    `json:"seller_id" db:"seller_id"`
	DropshipperId int    `json:"dropshipper_id" db:"dropshipper_id"`
	OrderDetails  string `json:"order_details" db:"order_details"`
	Region        string `json:"region" db:"region"`
}

type ParamsCreateDelivery struct {
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
	DeliveryPaymentMethod string  `json:"delivery_payment_method" form:"delivery_payment_method" db:"delivery_payment_method"`

	// We need new fields: Picture, and Item Description
	ImageUrl        string `json:"image_url" form:"image_url" db:"image_url" required:"false"`
	ItemDescription string `json:"item_description" form:"item_description" db:"item_description"`
}

type ParamsCreateParcel struct {
	DeliveryOption  string  `json:"delivery_option" form:"delivery_option" db:"delivery_option"`
	SellerId        int     `json:"seller_id" form:"seller_id" db:"seller_id"`
	DropshipperId   int     `json:"dropshipper_id" form:"dropshipper_id" db:"dropshipper_id"`
	Name            string  `json:"name" form:"name" db:"name"`
	ContactNumber   string  `json:"contact_number" form:"contact_number" db:"contact_number"`
	Address         string  `json:"address" form:"address" db:"address"`
	Note            string  `json:"note" form:"note" db:"note"`
	Region          string  `json:"region" form:"region" db:"region"`
	ServiceFee      float64 `json:"service_fee" form:"service_fee" db:"service_fee"`
	DeclaredAmount  float64 `json:"declared_amount" form:"declared_amount" db:"declared_amount"`
	DeliveryDetails string  `json:"delivery_details" form:"delivery_details" db:"delivery_details"`
}

type Top10Seller struct {
	SellerId   int    `json:"id" db:"id"`
	Seller     string `json:"name" db:"name"`
	Deliveries int    `json:"deliveries" db:"deliveries"`
	Sales      int    `json:"sales" db:"sales"`
}

type ValidateDropshipper struct {
	Id     int    `json:"id" db:"id"`
	Region string `json:"region" db:"region"`
}

type SearchTransactionFilter struct {
	UserId  int `json:"user_id" db:"user_id"`
	TranNum int `json:"tran_num" db:"tran_num"`
}

func (s *SearchTransactionFilter) Populate(c *fiber.Ctx) error {
	reqUserId := c.Query("user_id")
	reqTranNum := c.Query("tran_num")

	if reqUserId != "" {
		userId, err := strconv.Atoi(reqUserId)
		if err != nil {
			return err
		} else {
			s.UserId = userId
		}
	}
	if reqTranNum != "" {
		tranNum, err := strconv.Atoi(reqTranNum)
		if err != nil {
			return err
		} else {
			s.TranNum = tranNum
		}
	}

	return nil
}
