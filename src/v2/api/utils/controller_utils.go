package utils

import (
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Caller struct {
	UserId   int
	UserType string
}

func GetCallerDetails(c *fiber.Ctx) *Caller {
	return &Caller{
		UserId:   c.Locals("tokenExtractedUserId").(int),
		UserType: c.Locals("tokenExtractedUserType").(string),
	}
}

// RespondError - returns an error formatted JSON using Fiber.Ctx
func RespondError(c *fiber.Ctx, operationStatus string, apiErr *ApplicationError) error {
	r := response_builder.Response{
		HttpCode:        apiErr.HttpStatus,
		ResponseMessage: apiErr.Message,
		Data:            apiErr.Error,
		OperationStatus: operationStatus,
		Pagination:      nil,
	}
	r.SetErrors(apiErr.Error)
	return c.Status(apiErr.HttpStatus).JSON(r)
}

// Respond - returns an error formatted JSON using Fiber.Ctx
func Respond(c *fiber.Ctx, operationStatus string, responseMessage string, data interface{}) error {
	r := response_builder.Response{
		HttpCode:        http.StatusOK,
		ResponseMessage: responseMessage,
		Data:            data,
		OperationStatus: operationStatus,
	}
	return c.Status(http.StatusOK).JSON(r)
}


