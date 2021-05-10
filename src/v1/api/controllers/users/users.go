package users

import (
	"fmt"
	BankTypeModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/bank_type"
	RegionModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/region"
	UserTypeModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/user_type"
	"strconv"

	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	UserModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/users"
	UserService "github.com/dembygenesis/droppy-prulife/src/v1/api/services/users"

	// Service "github.com/dembygenesis/droppy-prulife/services/users"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	// "fmt"
	"github.com/gofiber/fiber/v2"
	// "github.com/dembygenesis/droppy-prulife/utilities/string"
	// "log"
)


func GetUserInfo(c *fiber.Ctx) error {

	u := UserModel.User{ID: c.Locals("tokenExtractedUserId").(int)}

	res, err := u.GetDetailsById()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to get the user details by id",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)

	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Successfully created the orders",
	}

	fmt.Println("res", res)

	r.Data = struct {
		UserInfo struct {
			Token       string                          `json:"token,omitempty"`
			UserDetails interface{} `json:"userDetails,omitempty"`
		} `json:"userInfo,omitempty"`
	}{UserInfo: struct {
		Token       string                          `json:"token,omitempty"`
		UserDetails interface{} `json:"userDetails,omitempty"`
	}{
		Token:       c.Get("authorization"),
		UserDetails: res,
	}}

	return c.JSON(r)
}

// Login attempts to authenticate your username and password
func Login(c *fiber.Ctx) error {

	var paramsLogin UserMiddleware.ParamsLogin

	// Ignore error handling because we already checked in the middleware
	_ = c.BodyParser(&paramsLogin)

	jwtToken, responseLoginUserInfo, err := UserService.Login(paramsLogin.Email, paramsLogin.Password)

	// Error
	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed login",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	// Success
	r := UserModel.ResponseUserDetailsDisplay{
		Response: response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Here's your credentials",
		},
		Data: struct {
			UserInfo struct {
				Token       string                          `json:"token,omitempty"`
				UserDetails UserModel.ResponseLoginUserInfo `json:"userDetails,omitempty"`
			} `json:"userInfo,omitempty"`
		}{UserInfo: struct {
			Token       string                          `json:"token,omitempty"`
			UserDetails UserModel.ResponseLoginUserInfo `json:"userDetails,omitempty"`
		}{
			Token:       jwtToken,
			UserDetails: responseLoginUserInfo,
		}},
	}

	return c.JSON(r)
}

// GetUserTypes returns all the user types
func GetUserTypes(c *fiber.Ctx) error {
	userType := UserTypeModel.UserType{}

	res, err := userType.GetAll()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to get the user types",
		}
		r.AddErrors("Something went wrong when trying to fetch the user types: " + err.Error())

		return c.JSON(r)

	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Here's the user types",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

// GetBankTypes returns all the banks
func GetBankTypes(c *fiber.Ctx) error {
	bankType := BankTypeModel.BankType{}

	res, err := bankType.GetAll()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to get the bank types",
		}
		r.AddErrors("Something went wrong when trying to fetch the bank types: " + err.Error())

		return c.JSON(r)

	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Here's the bank types",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

func GetRegions(c *fiber.Ctx) error {


	fmt.Println("Piss me off")

	// return c.SendString("Hello!")
	region := RegionModel.Region{}

	res, err := region.GetAll()


	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to get the regions",
		}
		r.AddErrors("Something went wrong when trying to fetch the regions: " + err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Here's the regions",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

// GetOne returns a user's detailed information
func GetOne(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))

	user := UserModel.User{ID: userId}

	userIdExtracted := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	if userType == "Dropshipper" || userType == "Seller" {
		// They can only edit themselves
		if userIdExtracted != userId {
			r := response_builder.Response{}
			r.HttpCode = 200
			r.ResponseMessage = "Something went wrong when trying to fetch the transactions"
			r.AddErrors("Dropshippers and Sellers can only edit themselves")

			return c.JSON(r)
		}
	}

	res, err := user.GetOne()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to fetch the user",
		}
		r.AddErrors("Something went wrong when trying to fetch the user " + err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Here's the user details",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

func GetAll2(c *fiber.Ctx) error {
	response := response_builder.Response{}
	response.HttpCode = 200

	u := UserModel.User{}

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

	fmt.Println("page", page, "rows", rows)

	// Attempt to fetch all users
	res, pagination, err := u.GetAll2(page, rows)

	if err != nil {
		response.ResponseMessage = "Something went wrong when trying to fetch the users"
		response.AddErrors(err.Error())

		return c.JSON(response)
	}

	response.Pagination =& pagination
	response.ResponseMessage = "Here's the users"

	userList := UserModel.ResponseUserList{
		Data:     res,
		Response: response,
	}

	return c.JSON(userList)
}

func GetAll(c *fiber.Ctx) {
	response := response_builder.Response{}
	response.HttpCode = 200

	// Attempt to fetch all users
	res, err := UserService.GetAll()

	if err != nil {
		response.ResponseMessage = "Something went wrong when trying to fetch the users"
		response.AddErrors(err.Error())

		c.JSON(response)
		return
	}

	response.ResponseMessage = "Here's the users"

	userList := UserModel.ResponseUserList{
		Data:     res,
		Response: response,
	}

	c.JSON(userList)
}
func GetAllByType(c *fiber.Ctx) error {
	response := response_builder.Response{}
	response.HttpCode = 200

	userType := c.Params("type")

	u := UserModel.User{}

	// Attempt to fetch all users
	res, err := u.GetAllByUserType(userType)

	fmt.Println("res", len(*res))

	if err != nil {
		response.ResponseMessage = "Something went wrong when trying to fetch the users by '" + userType + "'"
		response.AddErrors(err.Error())

		return c.JSON(response)
	}

	r := response_builder.Response{
		HttpCode:        200,
		ResponseMessage: "Here's the '" + userType + "'(s)",
		Data: make([]UserModel.UserListDisplay, 0),
	}

	if len(*res) == 0 {
		r.Data = make([]UserModel.UserListDisplay, 0)
	} else {
		r.Data = res
	}

	return c.JSON(r)
}

func Create(c *fiber.Ctx) error {
	var paramsInsert UserModel.ParamsInsert

	// Ignore error validation because we already did that in the middleware
	_ = c.BodyParser(&paramsInsert)

	// Attempt to insert the user
	paramsInsert.CreatedBy = c.Locals("tokenExtractedUserId").(int)

	user := UserModel.User{}

	_, err := user.Create(paramsInsert)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed Delete",
		}
		r.AddErrors("Something went wrong when trying to update the user " + err.Error())

		return c.JSON(r)

	}

	r := response_builder.ResponseSuccessOperation{
		Response: response_builder.Response{
			HttpCode:        200,
			OperationStatus: "INSERT_SUCCESS",
			ResponseMessage: "Successfully created the user",
		},
		Data: []string{"Successfully created the user"},
	}

	return c.JSON(r)
}

func Update(c *fiber.Ctx) error {
	var paramsUpdate UserModel.ParamsUpdate

	// Ignore error validation because we already did that in the middleware
	_ = c.BodyParser(&paramsUpdate)

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	if userType == "Dropshipper" || userType == "Seller" || userType == "Rider" {
		// They can only edit themselves
		if userId != paramsUpdate.ID {
			r := response_builder.Response{}
			r.HttpCode = 200
			r.ResponseMessage = "Something went wrong when trying to fetch the transactions"
			r.AddErrors("Dropshippers, Sellers, and Riders can only edit themselves")

			return c.JSON(r)
		}
	}

	_, err := paramsUpdate.Update(userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed Update",
		}
		r.AddErrors("Something went wrong when trying to update the user " + err.Error())

		return c.JSON(r)
	}

	r := response_builder.ResponseSuccessOperation{
		Response: response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Successfully updated the user",
		},
		Data: []string{"Successfully updated the user"},
	}
	r.SetUpdateSuccess()

	return c.JSON(r)
}

func Delete(c *fiber.Ctx) error {

	var paramsDelete UserModel.ParamsDelete

	// Ignoring errors because assumed ok in the middleware
	_ = c.BodyParser(&paramsDelete)

	// Attempt to delete (void) user
	user := UserModel.User{ID: paramsDelete.ID}

	_, err := user.Delete()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed Delete",
		}
		r.AddErrors("Something went wrong when trying to delete the user " + err.Error())

		return c.JSON(r)
	}

	r := response_builder.ResponseSuccessOperation{
		Response: response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Successfully deleted the user",
		},
		Data: []string{"Successfully deleted the user"},
	}

	return c.JSON(r)
}
