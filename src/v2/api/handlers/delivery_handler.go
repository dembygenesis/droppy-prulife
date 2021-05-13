package handlers

import (
	"github.com/dembygenesis/droppy-prulife/src/v2/api/config"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/services"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"github.com/dembygenesis/droppy-prulife/utilities/aws/s3"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type DeliveryHandler interface {
	// Implement controller method signatures
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type deliveryHandler struct {
	service services.Service
}

func NewDeliveryHandler(service services.Service) DeliveryHandler {
	return &deliveryHandler{service}
}

func (h *deliveryHandler) Update(c *fiber.Ctx) error {
	var appError *utils.ApplicationError

	// Validate body
	var body delivery.RequestUpdateDelivery
	err := c.BodyParser(&body)
	if err != nil {
		return utils.RespondError(c, "UPDATE_FAILED", &utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}
	caller := utils.GetCallerDetails(c)
	body.CreatedByUserType = caller.UserType

	appError = h.service.Update(&body)
	if appError != nil {
		return utils.RespondError(c, config.InsertFailed, appError)
	}

	return utils.Respond(c, config.UpdateSuccess, "Successfully created the delivery", nil)

	return nil
}

func (h *deliveryHandler) Create(c *fiber.Ctx) error {
	var appError *utils.ApplicationError

	// Validate body
	var body delivery.RequestCreateDelivery
	err := c.BodyParser(&body)
	if err != nil {
		return utils.RespondError(c, "INSERT_FAILED", &utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}
	caller := utils.GetCallerDetails(c)
	body.CreatedByUserType = caller.UserType

	// Validate file format to be image
	file, err := c.FormFile("image")
	if err != nil {
		appError = &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    "image is invalid",
			Error:      err,
		}
		return utils.RespondError(c, config.InsertFailed, appError)
	}
	// Decode the image file
	err = s3.DecodeImage(file)
	if err != nil {
		appError = &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    "unable to decode the image",
			Error:      err,
		}
		return utils.RespondError(c, config.InsertFailed, appError)
	}
	// Insert
	appError = h.service.Create(&body, file)
	if appError != nil {
		return utils.RespondError(c, config.InsertFailed, appError)
	}

	return utils.Respond(c, config.InsertSuccess, "Successfully created the delivery", nil)
}