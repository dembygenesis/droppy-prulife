package app

import (
	CoinTransactionRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/coin_transactions"
	DeliveryRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/deliveries"
	OrderRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/orders"
	ProductRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/products"
	SysParamRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/sysparam"
	TransactionRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/transactions"
	UserRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/users"
	WithdrawalRoutes "github.com/dembygenesis/droppy-prulife/src/v1/api/routes/withdrawals"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func mapUrlsV1(app *fiber.App) {
	// Initial route group
	api := app.Group("/api/v1", cors.New(), logger.New())

	// Routes
	UserRoutes.BindRoutes(api)
	ProductRoutes.BindRoutes(api)
	CoinTransactionRoutes.BindRoutes(api)
	TransactionRoutes.BindRoutes(api)
	OrderRoutes.BindRoutes(api)
	DeliveryRoutes.BindRoutes(api)
	WithdrawalRoutes.BindRoutes(api)
	SysParamRoutes.BindRoutes(api)
}