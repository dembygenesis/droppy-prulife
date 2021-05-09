package sysparam

import (
	SysParamController "github.com/dembygenesis/droppy-prulife/src/v1/api/controllers/sysparam"
	UserMiddleware "github.com/dembygenesis/droppy-prulife/src/v1/api/middlewares/users"
	"github.com/gofiber/fiber/v2"
)

func BindRoutes(api fiber.Router) {
	api = api.Group("/sysparam")

	api.Get("/", UserMiddleware.RoleMiddlewareV2([]string{"Admin"}), SysParamController.GetAll)
	api.Put("/", UserMiddleware.RoleMiddlewareV2([]string{"Admin"}), SysParamController.Update)
}