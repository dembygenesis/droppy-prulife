package withdrawals

import (
	"encoding/json"
	WithdrawalModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/withdrawals"
	ResponseBuilder "github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Update(c *fiber.Ctx) error {
	var paramsUpdateWithdrawal WithdrawalModel.ParamsUpdateWithdrawal

	c.BodyParser(&paramsUpdateWithdrawal)

	var paramErrors []string

	if paramsUpdateWithdrawal.ID == 0 {
		paramErrors = append(paramErrors, "id empty")
	}

	if paramsUpdateWithdrawal.Status == "" {
		paramErrors = append(paramErrors, "status empty")
	}

	if paramsUpdateWithdrawal.Status == "Voided" && paramsUpdateWithdrawal.VoidOrRejectReason == "" {
		paramErrors = append(paramErrors, "void_or_reject_reason empty")
	}

	if paramsUpdateWithdrawal.Status == "Completed" && paramsUpdateWithdrawal.ReferenceNumber == "" {
		paramErrors = append(paramErrors, "reference_number empty")
	}

	if len(paramErrors) > 0 {
		r := ResponseBuilder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your withdrawal",
		}
		r.SetErrors(paramErrors)

		return c.JSON(r)
	}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	w := WithdrawalModel.Withdrawal{UserID: userId, UserType: userType}

	res, err := w.Update(paramsUpdateWithdrawal)

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to update the withdrawal"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully Updated Your Withdrawal!"
	r.OperationStatus = "UPDATE_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}

func Create(c *fiber.Ctx) error {
	var paramsCreateWithdrawal WithdrawalModel.ParamsCreateWithdrawal

	c.BodyParser(&paramsCreateWithdrawal)

	var paramErrors []string

	if paramsCreateWithdrawal.Amount < 1 {
		paramErrors = append(paramErrors, "amount empty")
	}

	if paramsCreateWithdrawal.BankNo == "" {
		paramErrors = append(paramErrors, "bank_no empty")
	}

	if paramsCreateWithdrawal.BankTypeId == 0 {
		paramErrors = append(paramErrors, "bank_type_id empty")
	}

	if paramsCreateWithdrawal.BankAccountName == "" {
		paramErrors = append(paramErrors, "bank_account_name empty")
	}

	if len(paramErrors) > 0 {
		r := ResponseBuilder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to create your withdrawal",
		}
		r.SetErrors(paramErrors)

		return c.JSON(r)
	}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	w := WithdrawalModel.Withdrawal{UserID: userId, UserType: userType}

	res, err := w.Create(paramsCreateWithdrawal)

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to add the withdrawal"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully Created Your Withdrawal!"
	r.OperationStatus = "INSERT_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}

func GetAll(c *fiber.Ctx) error {

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
	Filters
	*/
	var filters []string

	f := c.Query("search")

	if f != "" {
		err := json.Unmarshal([]byte(f), &filters)

		if err != nil {
			r := ResponseBuilder.Response{}
			r.HttpCode = 200
			r.ResponseMessage = "Something went wrong when trying to decode the filters"
			r.AddErrors("Something went wrong when trying to decode the filters " + err.Error() + "with value of : " + f)

			return c.JSON(r)
		}
	}

	w := WithdrawalModel.Withdrawal{UserID: userId, UserType: userType}

	res, pagination, err := w.GetAll(page, rows, filters)

	if err != nil {
		r := ResponseBuilder.Response{}
		r.HttpCode = 200
		r.ResponseMessage = "Something went wrong when trying to fetch the transactions"
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{
		HttpCode:        200,
		ResponseMessage: "Here's the withdrawals",
		Pagination:      &pagination,
	}

	if len(*res) == 0 {
		r.Data = make([]WithdrawalModel.ResponseWithdrawalList, 0)
	} else {
		r.Data = res
	}

	return c.JSON(r)
}