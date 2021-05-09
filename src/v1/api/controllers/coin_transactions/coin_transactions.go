package coin_transactions

import (
	"fmt"
	CoinTransactionsModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/coin_transactions"
	ResponseBuilder "github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
)

func GetAll(c *fiber.Ctx) error {
	t := CoinTransactionsModel.CoinTransaction{}

	// Attempt to fetch all transactions
	res, err := t.GetAll()

	fmt.Println(res)

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to fetch the coin transactions"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's the coin transactions"
	r.SetResponseData(res)

	return c.JSON(r)
}

func Create(c *fiber.Ctx) error {

	return nil
	/*transaction := CoinTransactionsModel.Transaction{}

	res, err := transaction.Create()

	if err != nil {
		r := ResponseBuilder.Response{
			HttpCode: 200,
			ResponseMessage: "Something went wrong when trying to insert a transaction",
		}
		r.AddErrors(err.Error())

		c.JSON(r)

		return
	}

	r := ResponseBuilder.Response{
		HttpCode: 200,
		ResponseMessage: "Successfully added a new transaction",
	}
	r.SetResponseData(res)

	c.JSON(r)*/
}