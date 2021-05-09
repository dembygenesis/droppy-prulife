package orders

import (
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
)
import OrdersModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/orders"

func Create(c *fiber.Ctx) error {

	var paramsOrder OrdersModel.ParamsOrder

	// Map body
	err := c.BodyParser(&paramsOrder)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate add orders",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	var errors []string

	// Validate params
	if paramsOrder.UserId == 0 {
		errors = append(errors, "user_id must not be empty")

	}

	if paramsOrder.OrderDetails == "" {
		errors = append(errors, "order_details must not be empty")
	}

	if len(errors) > 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate add orders",
		}
		r.SetErrors(errors)
		return c.JSON(r)
	}

	return c.Next()
}