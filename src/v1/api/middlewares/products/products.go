package products

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)
import ProductModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/products"
import ProductTypeModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/product_types"
import "github.com/dembygenesis/droppy-prulife/utilities/response_builder"

func DeleteMiddleware(c *fiber.Ctx) error {
	var paramsDelete ProductModel.ParamsDelete

	// Attempt to parse args
	err := c.BodyParser(&paramsDelete)

	fmt.Println(":paramsDelete", paramsDelete)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate delete product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	// Validate ID
	product := ProductModel.Product{ID: paramsDelete.ID}

	res, err := product.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}
		r.AddErrors("id must be valid")


		return c.JSON(r)
	}

	return c.Next()
}


func UpdateMiddleware(c *fiber.Ctx) error {
	var paramsUpdate ProductModel.ParamsUpdate

	// Parse Body
	err := c.BodyParser(&paramsUpdate)

	fmt.Println(":paramsUpdate", paramsUpdate)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate update product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	// Handle errors
	var emptyParameters []string

	// Parameters must not be empty
	if paramsUpdate.ID == 0 {
		emptyParameters = append(emptyParameters, "id must not be empty")
	}

	if paramsUpdate.Name == "" {
		emptyParameters = append(emptyParameters, "name must not be empty")
	}

	fmt.Println("paramsUpdate.ProductTypeId", paramsUpdate.ProductTypeId)

	if paramsUpdate.ProductTypeId == 0 {
		emptyParameters = append(emptyParameters, "product_type_id must not be empty")
	}

	if len(emptyParameters) > 0 {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	// Validate Product Type
	productType := ProductTypeModel.ProductType{ID: paramsUpdate.ProductTypeId}

	res, err := productType.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product type validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product type validation failed",
		}
		r.AddErrors("product_type_id must be valid")


		return c.JSON(r)
	}

	product := ProductModel.Product{ID: paramsUpdate.ID, Name: paramsUpdate.Name}

	// Validate ID
	res, err = product.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}
		r.AddErrors("id must be valid")

		return c.JSON(r)
	}

	// Validate name
	res, err = product.UniqueNameExceptOwn()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}
		r.AddErrors("name must be unique")


		return c.JSON(r)

	}

	return c.Next()
}

func CreateMiddleware(c *fiber.Ctx) error {

	var paramsCreate ProductModel.ParamsCreate

	err := c.BodyParser(&paramsCreate)

	if err != nil {
		r := response_builder.Response{
			HttpCode:        200,
			ResponseMessage: "Failed to validate create product",
		}
		r.AddErrors(err.Error())

		return c.JSON(r)
	}

	var emptyParameters []string

	// Parameters must not be empty
	if paramsCreate.Name == "" {
		emptyParameters = append(emptyParameters, "name must not be empty")
	}

	if paramsCreate.ProductTypeId == 0 {
		emptyParameters = append(emptyParameters, "product_type_id must not be empty")
	}

	if len(emptyParameters) > 0 {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	// Validate Product Type
	productType := ProductTypeModel.ProductType{ID: paramsCreate.ProductTypeId}

	res, err := productType.ValidId()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product type validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product type validation failed",
		}
		r.AddErrors("product_type_id must be valid")


		return c.JSON(r)
	}

	// Validate name
	product := ProductModel.Product{Name: paramsCreate.Name}
	res, err = product.UniqueName()

	if err != nil {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}

		for _,v := range emptyParameters {
			r.AddErrors(v)
		}

		return c.JSON(r)
	}

	if res == false {
		r := response_builder.Response{
			HttpCode: 200,
			ResponseMessage: "Product validation failed",
		}
		r.AddErrors("name must be unique")


		return c.JSON(r)
	}

	return c.Next()
}