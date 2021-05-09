package withdrawals

import (
	WithdrawalController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/withdrawals"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/withdrawal")

	api.Get("/", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin"}), WithdrawalController.GetAll)
	api.Post("/", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin"}), WithdrawalController.Create)
	api.Put("/", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin"}), WithdrawalController.Update)
}
