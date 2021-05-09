package users

import (
	"fmt"
	BankTypeModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/bank_type"
	RegionModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/region"
	UserTypeModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/user_type"
	UserModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/users"
	"github.com/dembygenesis/droppy-prulife/utilities/response_builder"
	StringUtility "github.com/dembygenesis/droppy-prulife/utilities/string"
	"strconv"
	"time"

	//	"fmt"
	"github.com/gofiber/fiber/v2"
)

func CreateMiddleware(c *fiber.Ctx) error {

	var paramsInsert UserModel.ParamsInsert

	// Make sure all parameters are present
	err := c.BodyParser(&paramsInsert)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "something went wrong when trying to parse the update parameters",
		}
		r.AddErrors("something went wrong when trying to parse the update parameters: " + err.Error())

		return c.JSON(r)

	}

	// Make sure there are no empty params
	emptyFields := paramsInsert.NoEmptyFields()

	if len(emptyFields) != 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}

		for _, val := range emptyFields {
			r.AddErrors(val)
		}

		return c.JSON(r)
	}

	// User Type Id Must Be Valid
	userType := UserTypeModel.UserType{ID: paramsInsert.UserTypeId}

	res, err := userType.ValidID()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to check the user_type_id")

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("user_type_id id must be valid")

		return c.JSON(r)
	}

	// Region ID must be valid
	requestedUserType, err := userType.GetNameById()

	if requestedUserType == "Rider" {

		fmt.Println("paramsInsert.RegionId", paramsInsert.RegionId)

		// Require Region
		if paramsInsert.RegionId == 9999 || paramsInsert.RegionId == 0 {
			r := response_builder.Response{
				HttpCode:        200,
				ResponseMessage: "Failed user create",
			}
			r.AddErrors("region_id is missing")

			return c.JSON(r)
		}

		// Validate region
		region := RegionModel.Region{ID: paramsInsert.RegionId}

		res, err = region.ValidID()

		// Region ID database error
		if err != nil {
			r := response_builder.Response{
				HttpCode:        200,
				ResponseMessage: "Failed user create",
			}
			r.AddErrors("something went wrong when trying to validate the region_id")

			return c.JSON(r)
		}

		// Region ID not found
		if res == false {
			r := response_builder.Response{
				HttpCode:        200,
				ResponseMessage: "Failed user create",
			}
			r.AddErrors("region_id id must be valid")

			return c.JSON(r)
		}
	}

	// Bank Type Id Must Be Valid
	bankType := BankTypeModel.BankType{ID: paramsInsert.BankTypeId}

	res, err = bankType.ValidID()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to check the bank_type_id")

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("bank_type_id id must be valid")

		return c.JSON(r)
	}

	// Birthday must be a valid date format
	_, err = time.Parse("2006-01-02", paramsInsert.Birthday)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to parse the birthday: " + err.Error())

		return c.JSON(r)
	}

	// Gender must be M or F
	validGender := []string{"M", "F"}

	stringInSlice := StringUtility.StringInSlice(paramsInsert.Gender, validGender)

	if stringInSlice == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("Gender must be M or F")

		return c.JSON(r)
	}

	// Make sure email is not taken
	user := UserModel.User{Email: paramsInsert.Email}

	res, err = user.EmailNotTaken()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("something went wrong when trying to validate the email: " + err.Error())

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user create",
		}
		r.AddErrors("email already taken")

		return c.JSON(r)
	}

	// Validate region only IF user type is "Rider"
	/*requestedUserType, err := userType.GetNameById()

	if requestedUserType == "Rider" {
		// Require Region
		if paramsInsert.RegionId == 9999 {

		}
	}*/

	return c.Next()
}

func UpdateMiddleware(c *fiber.Ctx) error {

	var paramsUpdate UserModel.ParamsUpdate

	// Make sure all parameters are present
	err := c.BodyParser(&paramsUpdate)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "something went wrong when trying to parse the update parameters",
		}
		r.AddErrors("something went wrong when trying to parse the update parameters: " + err.Error())

		return c.JSON(r)
	}

	// Make sure there are no empty params
	emptyFields := paramsUpdate.NoEmptyFields()

	if len(emptyFields) != 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}

		for _, val := range emptyFields {
			r.AddErrors(val)
		}

		return c.JSON(r)
	}

	// User id must be valid
	user := UserModel.User{ID: paramsUpdate.ID}

	res, err := user.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}
		r.AddErrors("something went wrong when trying to check the user")

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed login",
		}
		r.AddErrors("user id must be valid")

		return c.JSON(r)
	}

	// User Type Id Must Be Valid
	userType := UserTypeModel.UserType{ID: paramsUpdate.UserTypeId}

	res, err = userType.ValidID()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}
		r.AddErrors("something went wrong when trying to check the user_type_id")

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed login",
		}
		r.AddErrors("user_type_id id must be valid")

		return c.JSON(r)
	}

	// Bank Type Id Must Be Valid
	bankType := BankTypeModel.BankType{ID: paramsUpdate.BankTypeId}

	res, err = bankType.ValidID()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}
		r.AddErrors("something went wrong when trying to check the bank_type_id")

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}
		r.AddErrors("bank_type_id id must be valid")

		return c.JSON(r)
	}

	// Birthday must be a valid date format
	_, err = time.Parse("2006-01-02", paramsUpdate.Birthday)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}
		r.AddErrors("something went wrong when trying to parse the birthday: " + err.Error())

		return c.JSON(r)
	}

	// Gender must be M or F
	validGender := []string{"M", "F"}

	stringInSlice := StringUtility.StringInSlice(paramsUpdate.Gender, validGender)

	if stringInSlice == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed user update",
		}
		r.AddErrors("Gender must be M or F")

		return c.JSON(r)

	}

	return c.Next()
}

func DeleteMiddleware(c *fiber.Ctx) error {

	var paramsDelete UserModel.ParamsDelete

	// Validate param "id"
	err := c.BodyParser(&paramsDelete)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed login",
		}
		r.AddErrors("something went wrong when trying to parse the user id")

		return c.JSON(r)
	}

	if paramsDelete.ID == 0 {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed login",
		}
		r.AddErrors("id parameter empty")

		return c.JSON(r)
	}

	// Validate ID
	user := UserModel.User{ID: paramsDelete.ID}

	exists, err := user.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed Delete",
		}
		r.AddErrors("something went wrong when trying to validate the user id")

		return c.JSON(r)

	}

	if exists == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed Delete",
		}
		r.AddErrors("user id doesn't exist")

		return c.JSON(r)
	}

	return c.Next()
}

// Try to add another object to params
func RoleMiddleware(role string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Attempt to validate token
		token := c.Get("authorization")

		user := UserModel.User{Token: token}

		userId, err := user.ValidateToken(role)

		if err != nil {
			var response = response_builder.Response{
				HttpCode:        401,
				ResponseMessage: "Unauthorizedd",
			}
			response.AddErrors(err.Error())

			return c.Status(401).JSON(response)
		}

		c.Locals("tokenExtractedUserId", userId)
		return c.Next()
	}
}

func RoleMiddlewareV2(roles []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Attempt to validate token
		token := c.Get("authorization")

		user := UserModel.User{Token: token}

		userId, userType, err := user.ValidateTokenV2(roles)

		if err != nil {
			var response = response_builder.Response{
				HttpCode:        401,
				ResponseMessage: "Unauthorized",
			}
			response.AddErrors(err.Error())

			return c.Status(401).JSON(response)

		}

		c.Locals("tokenExtractedUserId", userId)
		c.Locals("tokenExtractedUserType", userType)

		return c.Next()
	}
}


func LoginValidation(c *fiber.Ctx) error  {

	var paramsLogin ParamsLogin
	var response response_builder.Response

	err := c.BodyParser(&paramsLogin)

	// Parsing must be fine
	if err != nil {
		response.HttpCode = 200
		response.ResponseMessage = "Error"
		response.AddErrors("Something went wrong with parsing the arguments for Login")

		return c.JSON(response)
	}

	// No empty inputs
	if paramsLogin.Email == "" || paramsLogin.Password == "" {
		response.HttpCode = 200
		response.ResponseMessage = "Error"
		response.AddErrors("email param is required")
		response.AddErrors("password param is required")

		return c.JSON(response)
	}

	// Email must exist in record
	user := UserModel.UserLogin{Email: paramsLogin.Email}

	exists, err := user.ValidEmail()

	if err != nil {
		response.HttpCode = 200
		response.ResponseMessage = "Syntax Error"
		response.AddErrors(err.Error())

		return c.JSON(response)
	}

	if exists == false {
		response.HttpCode = 200
		response.ResponseMessage = "Syntax Error"
		response.AddErrors("Email does not exist")

		return c.JSON(response)
	}

	return c.Next()
}

func ValidUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed validation",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	user := UserModel.User{ID: userId}
	res, err := user.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed validation",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed validation",
		}
		r.AddErrors("user does not exist")

		return c.JSON(r)
	}

	return c.Next()
}