package sysparam

import (
	SysParamModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/sysparam"
	DatabaseService "github.com/dembygenesis/droppy-prulife/utilities/database"
	ResponseBuilder "github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	"github.com/gofiber/fiber/v2"
)

func Update(c *fiber.Ctx) error {
	var updateParams SysParamModel.ParamsUpdateSysParam

	// Parse body params
	err := c.BodyParser(&updateParams)

	if err != nil {
		r := ResponseBuilder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to update the SysParam",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	// Validate entry
	hasEntry, err := DatabaseService.ValidEntry(updateParams.Key, "`key`", "sysparam")

	if err != nil || hasEntry == false {
		r := ResponseBuilder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to update the SysParam",
		}

		if err != nil {
			r.AddErrors(err.Error())
		}

		if hasEntry == false {
			r.AddErrors("No valid key found")
		}

		return c.JSON(r)
	}


	// Execute update
	s := SysParamModel.SysParam{
		Key:   updateParams.Key,
	}

	res, err := s.Update(updateParams.Value)

	if err != nil {
		r := ResponseBuilder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to update the SysParam",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Successfully updated the config parameters!"
	r.OperationStatus = "UPDATE_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}

func GetAll(c *fiber.Ctx) error {
	s := SysParamModel.SysParam{}

	res, err := s.GetAll()

	if err != nil {
		r := ResponseBuilder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the config parameters",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := ResponseBuilder.Response{}
	r.HttpCode = 200
	r.ResponseMessage = "Here's the config parameterse!"
	r.OperationStatus = "FETCH_SUCCESS"
	r.SetResponseData(res)

	return c.JSON(r)
}