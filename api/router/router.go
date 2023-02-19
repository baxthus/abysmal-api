package router

import (
	"github.com/Abysm0xC/abysmal-api/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Initialize() *fiber.App {
	app := fiber.New()

	app.Use(
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		}),
		compress.New(),
		recover.New(),
	)

	app.Static("/", "/public")

	app.Mount("/", handlers.Routes())

	return app
}
