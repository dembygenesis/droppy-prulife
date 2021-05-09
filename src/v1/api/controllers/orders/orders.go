package orders

import (
	OrdersModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/orders"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAll(c *fiber.Ctx) error {
	return nil
}

func Get(c *fiber.Ctx) error {
	var paramsGetOrderDetails OrdersModel.ParamsGetOrderDetails

	c.BodyParser(&paramsGetOrderDetails)

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)
	orderId, _ := strconv.Atoi(c.Params("id"))

	o := OrdersModel.Order{}

	res, err := o.GetDetails(orderId, userId, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Successfully queried the order",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

/*func GetAll(c *fiber.Ctx) {
	var paramsOrder OrdersModel.ParamsOrder

	c.BodyParser(&paramsOrder)

	paramsOrder.AdminId = c.Locals("tokenExtractedUserId").(int)

	o := OrdersModel.Order{}

	res, err := o.GetAll()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the orders",
		}
		r.AddErrors(err.Error())

		c.JSON(r)
		return
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Successfully queried the orders",
	}
	r.SetResponseData(res)

	c.JSON(r)
	return
}*/

func Create(c *fiber.Ctx) error {

	var paramsOrder OrdersModel.ParamsOrder

	c.BodyParser(&paramsOrder)

	paramsOrder.AdminId = c.Locals("tokenExtractedUserId").(int)

	o := OrdersModel.Order{}

	res, err := o.Create(&paramsOrder)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate add orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Successfully created the orders",
		OperationStatus: "INSERT_SUCCESS",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

func Update(c *fiber.Ctx) error {

	var paramsOrderUpdate OrdersModel.ParamsOrderUpdate

	c.BodyParser(&paramsOrderUpdate)

	paramsOrderUpdate.AdminId = c.Locals("tokenExtractedUserId").(int)

	var paramErrors []string

	if paramsOrderUpdate.OrderId == 0 {
		paramErrors = append(paramErrors, "order_id empty")
	}

	if paramsOrderUpdate.Description == "" {
		paramErrors = append(paramErrors, "void_or_reject_reason empty")
	}

	if len(paramErrors) > 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to update your order",
		}
		r.SetErrors(paramErrors)

		return c.JSON(r)
	}

	o := OrdersModel.Order{}

	res, err := o.Update(&paramsOrderUpdate)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate add orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Successfully voided the order",
		OperationStatus: "DELETE_SUCCESS",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}