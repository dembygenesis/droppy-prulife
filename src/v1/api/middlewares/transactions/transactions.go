package transactions

import (
	"fmt"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
)
import TransactionModels "github.com/dembygenesis/droppy-prulife/src/v1/api/models/transactions"

func Delete(c *fiber.Ctx) error {
	// var paramsTransactionDelete TransactionModels.ParamsTransactionDelete

	return c.Next()
}

func Create(c *fiber.Ctx) error {
	var paramsTransaction TransactionModels.ParamsTransaction

	// Attempt to parse parameters
	err := c.BodyParser(&paramsTransaction)

	fmt.Println(":paramsTransaction", paramsTransaction)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to parse the create transaction parameters",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	// All fields must not be empty (further validation will be in the database)
	fmt.Println(paramsTransaction)

	var errors []string

	if paramsTransaction.Amount < 1 {
		errors = append(errors, "amount must be greater than 1 ")
	}

	if paramsTransaction.CoinAmount < 1 {
		errors = append(errors, "coin amount must be greater than 1 ")
	}

	if paramsTransaction.AdminId == 0 {
		errors = append(errors, "admin_id must be be empty ")
	}

	if paramsTransaction.UserId == 0 {
		errors = append(errors, "user_id must be be empty ")
	}

	// Lol can't do anything about money in
	if paramsTransaction.BankTypeId == 0 {
		errors = append(errors, "bank_type_id must be be empty ")
	}

	if paramsTransaction.ReferenceNumber == "" {
		errors = append(errors, "reference_number must be be empty ")
	}

	if paramsTransaction.Description == "" {
		errors = append(errors, "description must be be empty ")
	}

	if len(errors) != 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate the create transaction params",
		}

		for _, v := range errors {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	return c.Next()
}