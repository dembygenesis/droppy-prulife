package deliveries

import (
	"bytes"
	"fmt"
	DeliveryModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/deliveries"
	"github.com/dembygenesis/droppy-prulife/utilities/aws/s3"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	ResponseBuilder "github.com/dembygenesis/droppy-prulife/utilities/response_builder"
)

func ServiceFee(c *fiber.Ctx) error {
	deliveryId, _ := strconv.Atoi(c.Params("id"))

	u := DeliveryModel.Delivery{ID: deliveryId}

	orderDetails := c.Query("order_details")

	res, err := u.GetServiceFee(orderDetails)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the service fee",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your service fee"
	r.SetResponseData(res)

	return c.JSON(r)
}

func Get(c *fiber.Ctx) error {
	deliveryId, _ := strconv.Atoi(c.Params("id"))

	u := DeliveryModel.Delivery{ID: deliveryId}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	res, err := u.GetDeliveryDetails(userId, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the delivery details",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your delivery details"
	r.SetResponseData(res)

	return c.JSON(r)
}

func Top10Sellers(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	res, err := u.GetDashboardDeliveryStatus(userId, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your store details"
	r.SetResponseData(res)

	return c.JSON(r)
}

func MyStore(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	res, err := u.GetDashboardDeliveryStatus(userId, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your store details"
	r.SetResponseData(res)

	return c.JSON(r)
}

func UpdateDelivery(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	var paramsUpdateDelivery DeliveryModel.ParamsUpdateDelivery

	err := c.BodyParser(&paramsUpdateDelivery)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your parcel",
		}
		r.AddErrors("Something went wrong when parsing the body for UpdateDelivery")
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	// Guard against param errors the GO way lol
	var paramErrors []string

	/**
	DeliveryId
	DeliveryStatus
	TrackingNumber
	 */

	if paramsUpdateDelivery.DeliveryId == 0 {
		paramErrors = append(paramErrors, "delivery_id empty")
	}

	if paramsUpdateDelivery.DeliveryStatus == "" {
		paramErrors = append(paramErrors, "delivery_status empty")
	}

	// Nah ignore this, errors will be caught by the backend anyway
	if paramsUpdateDelivery.TrackingNumber == "" {
		// paramErrors = append(paramErrors, "declared_amount empty")
	}

	if len(paramErrors) > 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to update your delivery",
		}
		r.SetErrors(paramErrors)

		return c.JSON(r)
	}

	userId := c.Locals("tokenExtractedUserId").(int)

	u.UserId = userId

	_, err = u.UpdateDelivery(paramsUpdateDelivery)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your parcel",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Successfully updated the delivery",
		OperationStatus: "UPDATE_SUCCESS",
	}

	return c.JSON(r)
}

func CreateOld(c *fiber.Ctx) error {

	var paramsCreateParcel DeliveryModel.ParamsCreateDelivery

	// Parse parameters
	err := c.BodyParser(&paramsCreateParcel)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	paramsCreateParcel.SellerId = c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	// Validate file format to be image
	file, err := c.FormFile("image")

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	// FUCK
	err = s3.IsMultipartImage(file)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	// FUCK

	openedFile, err := file.Open()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	_, _, err = image.Decode(openedFile)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	d := DeliveryModel.Delivery{}

	// Ensure no empty parameters
	missingParameters := d.ValidateCreate(paramsCreateParcel)

	if len(missingParameters) > 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(missingParameters)

		return c.JSON(r)
	}

	// Attempt to create
	err = d.Create(&paramsCreateParcel, file, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	return c.JSON(paramsCreateParcel)
}



func Create(c *fiber.Ctx) error {

	var paramsCreateParcel DeliveryModel.ParamsCreateDelivery

	// Parse parameters
	err := c.BodyParser(&paramsCreateParcel)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	paramsCreateParcel.SellerId = c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	// Validate file format to be image
	file, err := c.FormFile("image")

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(err)

		return c.JSON(r)
	}

	// Validate image file type
	err = s3.IsMultipartImage(file)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	d := DeliveryModel.Delivery{}

	// Ensure no empty parameters
	missingParameters := d.ValidateCreate(paramsCreateParcel)

	if len(missingParameters) > 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.SetErrors(missingParameters)

		return c.JSON(r)
	}

	// Attempt to create
	err = d.Create(&paramsCreateParcel, file, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your delivery",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully Created Your Delivery"
	r.OperationStatus = "INSERT_SUCCESS"
	r.SetResponseData(nil)

	return c.JSON(r)
}

func OrderParcel2(c *fiber.Ctx) error {

	// Make a service that will validate inputs lololol
	var paramsCreateParcel DeliveryModel.ParamsCreateDelivery

	// Parse file
	file, err := c.FormFile("image")

	if err != nil {
		return c.JSON(err)
	}

	// Parse multipart file
	multiPartFile, err := file.Open()

	if err != nil {
		return c.JSON(err)
	}

	// Convert to buffer
	size := file.Size
	buffer := make([]byte, size)
	multiPartFile.Read(buffer)

	conversion := bytes.NewReader(buffer)

	// s3.MultiPart(conversion)

	s3.UploadObjectBuffer("fuck.txt", conversion, "droppy.biz")

	return c.JSON(paramsCreateParcel)
}

func OrderParcel(c *fiber.Ctx) error {

	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)

	var paramsCreateParcel DeliveryModel.ParamsCreateParcel

	err := c.BodyParser(&paramsCreateParcel)

	/**
	Handle image
	 */

	// Handle image
	_, err = c.FormFile("image")
	sword, err := c.FormFile("dropshipper_id")

	// fmt.Println("file", file)
	fmt.Println("err", err)

	fmt.Println("paramsCreateParcel", paramsCreateParcel)
	fmt.Println("paramsCreateParcel.DropshipperId", paramsCreateParcel.DropshipperId)

	fmt.Println("sword", sword)

	return c.SendString("hello")

	/**
	end Handle image
	*/


	// Guard against param errors the GO way lol
	var paramErrors []string

	if paramsCreateParcel.DropshipperId == 0 {
		paramErrors = append(paramErrors, "dropshipper_id empty")
	}

	if paramsCreateParcel.ServiceFee == 0 {
		paramErrors = append(paramErrors, "service_fee empty")
	}

	if paramsCreateParcel.DeclaredAmount == 0 {
		paramErrors = append(paramErrors, "declared_amount empty")
	}

	if paramsCreateParcel.DeliveryDetails == "" {
		paramErrors = append(paramErrors, "delivery_details empty")
	}

	if paramsCreateParcel.ContactNumber == "" {
		paramErrors = append(paramErrors, "contact_number empty")
	}

	if paramsCreateParcel.Note == "" {
		paramErrors = append(paramErrors, "note empty")
	}

	fmt.Println(paramsCreateParcel)

	if len(paramErrors) > 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your parcel",
		}
		r.SetErrors(paramErrors)

		return c.JSON(r)
	}

	paramsCreateParcel.SellerId = userId

	res, err := u.CreateParcel(paramsCreateParcel)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your parcel",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully Created Your Parcel!"
	r.OperationStatus = "INSERT_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}

func OrderPackage(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)
	// userType := c.Locals("tokenExtractedUserType").(string)

	var paramsCreateOrder DeliveryModel.ParamsCreateOrder

	err := c.BodyParser(&paramsCreateOrder)

	paramsCreateOrder.SellerId = userId

	res, err := u.CreatePackage(paramsCreateOrder)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your order",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully Created Your Order!"
	r.OperationStatus = "INSERT_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}

func CoinTransactions(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	/**
	Pagination
	*/
	var page int
	var rows int

	// Set default rows to 100 if not paginated
	if c.Query("page") == "" {
		page = 0
	} else {
		page, _ = strconv.Atoi(c.Query("page"))
	}

	if c.Query("rows") == "" {
		rows = 100
	} else {
		rows, _ = strconv.Atoi(c.Query("rows"))

		if rows <= 0 {
			rows = 1000
		}
	}

	res, pagination, err := u.GetCoinTransactions(userId, userType, page, rows)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the coin transactions",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your coin transaction details"
	r.SetResponseData(res)
	r.Pagination = pagination

	return c.JSON(r)
}

func CoinTransactions2(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	/**
	Pagination
	*/
	var page int
	var rows int

	// Set default rows to 100 if not paginated
	if c.Query("page") == "" {
		page = 0
	} else {
		page, _ = strconv.Atoi(c.Query("page"))
	}

	if c.Query("rows") == "" {
		rows = 100
	} else {
		rows, _ = strconv.Atoi(c.Query("rows"))

		if rows <= 0 {
			rows = 1000
		}
	}

	res, pagination, err := u.GetCoinTransactions2(userId, userType, page, rows)

	fmt.Println("CoinTransactions err", err)

	if err != nil {
		fmt.Println("Here I AM")
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the coin transactions",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	} else {
		fmt.Println("Here I AM else", err)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your coin transaction details"
	r.SetResponseData(res)
	r.Pagination = pagination

	return c.JSON(r)
}

func Transactions(c *fiber.Ctx) error {
	u := DeliveryModel.Delivery{}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	/**
		Pagination
	 */
	var page int
	var rows int

	// Set default rows to 100 if not paginated
	if c.Query("page") == "" {
		page = 0
	} else {
		page, _ = strconv.Atoi(c.Query("page"))
	}

	if c.Query("rows") == "" {
		rows = 100
	} else {
		rows, _ = strconv.Atoi(c.Query("rows"))

		if rows <= 0 {
			rows = 1000
		}
	}

	/**
		Search
	*/
	var search string

	if c.Query("search") == "" {
		search = ""
	} else {
		search = c.Query("search")
	}

	deliveryStatus := c.Query("delivery_status")

	res, pagination, err := u.GetTransactions(userId, userType, deliveryStatus, search, page, rows)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)

	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's your store details"
	r.SetResponseData(res)
	r.Pagination = pagination

	return c.JSON(r)
}