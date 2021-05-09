package products

import (
	"fmt"
	ProductTypeModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/product_types"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"strings"
)
import ProductModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/products"
import "github.com/dembygenesis/droppy-prulife/utilities/response_builder"

func GetSellerList(c *fiber.Ctx) error {
	p := ProductModel.Product{}

	userId := c.Locals("tokenExtractedUserId").(int)

	res, err := p.GetSellerList(userId)

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Something went wrong when trying to fetch the Seller List",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Here's the Seller List",
		Data: res,
	}

	return c.JSON(r)
}

func GetInventory(c *fiber.Ctx) error {
	p := ProductModel.Product{}

	userId := c.Locals("tokenExtractedUserId").(int)
	userType := c.Locals("tokenExtractedUserType").(string)

	res, err := p.GetInventory(userId, userType)

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "something went wrong when trying to fetch the product inventory",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Here's the product inventory",
		Data: res,
	}

	return c.JSON(r)
}

func Delete(c *fiber.Ctx) error {
	var paramsDelete ProductModel.ParamsDelete

	// Parse params
	_ = c.BodyParser(&paramsDelete)

	product := ProductModel.Product{
		ID: paramsDelete.ID,
	}

	// Attempt to delete
	_, err := product.Delete()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "something went wrong when trying to delete the product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Delete Success!",
		Data: []string{"Successfully deleted the product"},
	}

	return c.JSON(r)
}

func Update(c *fiber.Ctx) error {
	var paramsUpdate ProductModel.ParamsUpdate

	_ = c.BodyParser(&paramsUpdate)

	product := ProductModel.Product{
		ID: paramsUpdate.ID,
		Name: paramsUpdate.Name,
		ProductTypeId: paramsUpdate.ProductTypeId,
		CreatedBy: c.Locals("tokenExtractedUserId").(int),
	}

	// Handle image
	file, err := c.FormFile("image")

	if err == nil {
		// Make folder if not existing.
		err = os.MkdirAll("./public/images", os.ModePerm)

		if err != nil {
			r := response_builder.Response{
				HttpCode: 200,
				ResponseMessage: "something went wrong when trying to make a public folder",
			}
			r.AddErrors(err.Error())

			return c.JSON(r)
		}

		fileDetails := strings.Split(file.Filename, ".")
		extension := fileDetails[len(fileDetails) - 1]
		newFileName := paramsUpdate.Name + "." + extension
		newFilePath := "./public/images/" + newFileName

		err = c.SaveFile(file, fmt.Sprintf("./%s", newFilePath))

		if err != nil {
			r := response_builder.Response{
				HttpCode: 200,
				ResponseMessage: "something went wrong when trying to save a file",
			}
			r.AddErrors(err.Error())

			return c.JSON(r)
		}

		// Proceed to update file directory silently
		product.Url = "/images/" + newFileName

		product.UpdateUrl()
	}

	_, err = product.Update()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "something went wrong when trying to update the product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Update Success!",
		Data: []string{"Successfully updated the product"},
	}

	return c.JSON(r)
}

func Create(c *fiber.Ctx) error {
	var paramsCreate ProductModel.ParamsCreate

	_ = c.BodyParser(&paramsCreate)

	product := ProductModel.Product{
		Name: paramsCreate.Name,
		ProductTypeId: paramsCreate.ProductTypeId,
		CreatedBy: c.Locals("tokenExtractedUserId").(int),
	}

	_, err := product.Create()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "something went wrong when trying to create the product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Create Success!",
		Data: []string{"Successfully added the product"},
	}

	return c.JSON(r)
}


func GetAllTypes(c *fiber.Ctx) error {
	p := ProductTypeModel.ProductType{}

	res, err := p.GetAll()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "failed to fetch all the product types",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Here's the product types",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

func GetOne(c *fiber.Ctx) error {

	productId, _ := strconv.Atoi(c.Params("id"))

	product := ProductModel.Product{ID: productId}

	res, err := product.GetOne()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "failed to update fetch the product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Here's the product",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

func GetAll(c *fiber.Ctx) error {
	product := ProductModel.Product{}

	res, err := product.GetAll(c.Query("filter"))

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "failed to update fetch all the products",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	r := response_builder.Response{
		HttpCode: 200,
		ResponseMessage: "Here's the products",
	}
	r.SetResponseData(res)

	return c.JSON(r)
}

