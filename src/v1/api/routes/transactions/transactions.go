package coin_transactions

import (
	TransactionController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/transactions"
	TransactionMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/transactions"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/transaction")

	// Crud
	api.Get("/", UserMiddleware.RoleMiddleware("Admin"), TransactionController.GetAll)
	api.Post("/", UserMiddleware.RoleMiddleware("Admin"), TransactionMiddleware.Create, TransactionController.Create)
	api.Delete("/", UserMiddleware.RoleMiddleware("Admin"), TransactionMiddleware.Delete, TransactionController.Delete)
}