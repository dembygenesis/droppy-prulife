package products

import (
	ProductController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/products"
	ProductMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/products"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/product")

	// Crud
	api.Get("/types", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Dropshipper", "Seller"}), ProductController.GetAllTypes)

	// Just lump the inventory here.
	api.Get("/inventory", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Dropshipper", "Seller"}), ProductController.GetInventory)
	api.Get("/seller-list", UserMiddleware.RoleMiddlewareV2([]string{"Dropshipper"}), ProductController.GetSellerList)

	api.Get("/", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Dropshipper", "Seller"}), ProductController.GetAll)
	api.Post("/", UserMiddleware.RoleMiddleware("Admin"), ProductMiddleware.CreateMiddleware, ProductController.Create)
	api.Put("/", UserMiddleware.RoleMiddleware("Admin"), ProductMiddleware.UpdateMiddleware, ProductController.Update)
	api.Delete("/", UserMiddleware.RoleMiddleware("Admin"), ProductMiddleware.DeleteMiddleware, ProductController.Delete)

	api.Get("/:id", UserMiddleware.RoleMiddleware("Admin"), ProductController.GetOne)
}