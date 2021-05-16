package app

import (
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/handlers"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/repositories/delivery/crud"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func mapUrlsV2(app *fiber.App) {

	deliveryHandler := handlers.NewDeliveryHandler(services.NewDeliveryService( crud.NewDeliveryRepository() ))

	api := app.Group("/api/v2", cors.New(), logger.New())

	api.Post("/delivery", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), deliveryHandler.Create)
	api.Put("/delivery", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin"}), deliveryHandler.Update)
}