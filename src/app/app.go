package app

import (
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"os"
)

func Start() {
	app := fiber.New(fiber.Config{
		BodyLimit: 20971520,
	})

	// Recover
	app.Use(recover.New())

	// Database
	database.EstablishConnection()

	// Routes
	mapUrlsV1(app)
	mapUrlsV2(app)

	// Public routes
	app.Static("/", "./public")
	app.Static("/prefix", "./public")
	app.Static("*", "./public/index.html")

	// Load port from env
	_ = godotenv.Load()

	// Boot app
	_ = app.Listen(":" + os.Getenv("PORT"))
}