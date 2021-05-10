package handlers

import (
	"github.com/dembygenesis/droppy-prulife/src/v2/api/config"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/services"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type DeliveryHandler interface {
	// Implement controller method signatures
	Create(c *fiber.Ctx) error
}

type deliveryHandler struct {
	service services.Service
}

func NewDeliveryHandler(service services.Service) DeliveryHandler {
	return &deliveryHandler{service}
}

func (h *deliveryHandler) Create(c *fiber.Ctx) error {
	var body delivery.CreateDelivery
	err := c.BodyParser(&body)
	if err != nil {
		return utils.RespondError(c, "INSERT_FAILED", &utils.ApplicationError{
			HttpStatus: http.StatusBadRequest,
			Message:    "bad_request",
			Error:      err,
		})
	}

	caller := utils.GetCallerDetails(c)
	body.SellerId = caller.UserId

	appError := h.service.Create(&body)


	if appError != nil {
		return utils.RespondError(c, config.InsertFailed, appError)
	}
	return utils.Respond(c, config.InsertSuccess, "Successfully created the delivery", nil)
}