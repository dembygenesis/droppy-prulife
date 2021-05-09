package coin_transactions

import (
	OrdersController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/orders"
	OrdersMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/orders"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/order")

	// Crud

	// api.Get("/", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller"}), OrdersController.GetAll)
	api.Post("/", UserMiddleware.RoleMiddleware("Admin"), OrdersMiddleware.Create, OrdersController.Create)
	api.Put("/", UserMiddleware.RoleMiddleware("Admin"), OrdersController.Update)

	api.Get("/", UserMiddleware.RoleMiddlewareV2([]string{"Admin"}), OrdersController.GetAll)
	api.Get("/:id", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin"}), OrdersController.Get)
}
