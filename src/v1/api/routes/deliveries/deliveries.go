package deliveries

import (
	DeliveriesController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/deliveries"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/delivery")

	api.Get("/my-store", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.MyStore)
	api.Get("/top-10-sellers", UserMiddleware.RoleMiddlewareV2([]string{"Admin"}), DeliveriesController.Top10Sellers)
	api.Get("/service-fee", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.ServiceFee)
	api.Get("/transactions", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin", "Rider"}), DeliveriesController.Transactions)
	api.Get("/coin-transactions", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.CoinTransactions)
	api.Get("/coin-transactions2", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.CoinTransactions2)
	api.Post("/order-package", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.OrderPackage)

	// Backup
	// api.Post("/order-parcel", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.OrderParcel)

	api.Post("/", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.Create)

	// api.Post("/order-parcel", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller"}), DeliveriesController.OrderParcel3)


	api.Put("/", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Admin"}), DeliveriesController.UpdateDelivery)

	api.Get("/:id", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper", "Seller", "Admin"}), DeliveriesController.Get)
}