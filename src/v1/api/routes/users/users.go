package users

import (
	UserController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/users"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {

	api = api.Group("/user")

	// Login routes
	api.Post("/login", UserMiddleware.LoginValidation, UserController.Login)
	api.Get("/refresh-data", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper", "Rider"}), UserController.GetUserInfo)

	// Crud
	api.Get("/user-type/:type", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper"}), UserController.GetAllByType)
	api.Get("/", UserMiddleware.RoleMiddleware("Admin"), UserController.GetAll2)
	// api.Get("/arsenic", UserMiddleware.RoleMiddleware("Admin"), UserController.GetAll2)
	api.Delete("/", UserMiddleware.RoleMiddleware("Admin"), UserMiddleware.DeleteMiddleware, UserController.Delete)

	// Update this to have Admin, Seller, and Dropshipper
	api.Put("/", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper", "Rider"}), UserMiddleware.UpdateMiddleware, UserController.Update)

	// Create User
	api.Post("/", UserMiddleware.RoleMiddleware("Admin"), UserMiddleware.CreateMiddleware, UserController.Create)

	// Options
	api.Get("/user-types", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper", "Rider"}), UserController.GetUserTypes)
	api.Get("/bank-types", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper", "Rider"}), UserController.GetBankTypes)
	api.Get("/regions", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper", "Rider"}), UserController.GetRegions)

	// Put this last because we might accidentally parse /user-types or /bank-types lol
	api.Get("/:id", UserMiddleware.RoleMiddlewareV2([]string{"Admin", "Seller", "Dropshipper", "Rider"}), UserMiddleware.ValidUser, UserController.GetOne)
}