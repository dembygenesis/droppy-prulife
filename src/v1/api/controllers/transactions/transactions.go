package transactions

import (
	"fmt"
	TransactionsModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/transactions"
	ResponseBuilder "github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
)

func GetAll(c *fiber.Ctx) error {
	t := TransactionsModel.Transaction{}

	// Attempt to fetch all transactions
	res, err := t.GetAll()

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to fetch the transactions"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's the transactions"
	r.SetResponseData(res)

	return c.JSON(r)
}

// Delete reverses a transaction and all it's underlying coin transactions.
func Delete(c *fiber.Ctx) error {
	var paramsTransactionDelete TransactionsModel.ParamsTransactionDelete

	c.BodyParser(&paramsTransactionDelete)

	// Overly optimistic
	t := TransactionsModel.Transaction{}

	// Attempt to fetch all transactions
	res, err := t.Delete(&paramsTransactionDelete)

	fmt.Println(res)

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to void the transactions"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully voided the transaction"
	r.OperationStatus = "INSERT_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}

// Create reverses a transaction and all it's underlying coin transactions.
func Create(c *fiber.Ctx) error {
	var paramsTransaction TransactionsModel.ParamsTransaction

	// Ignore error handling (done on middleware already)
	c.BodyParser(&paramsTransaction)

	// Attempt to insert
	t := TransactionsModel.Transaction{}

	res, err := t.Create(&paramsTransaction)

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to create a transactions"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Created the transaction"
	r.OperationStatus = "INSERT_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}