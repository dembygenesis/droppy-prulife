package coin_transactions

import (
	CoinTransactionController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/coin_transactions"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/coin-transaction")

	// Crud
	api.Get("/", UserMiddleware.RoleMiddleware("Admin"), CoinTransactionController.GetAll)
	api.Post("/", UserMiddleware.RoleMiddleware("Admin"), CoinTransactionController.Create)
}