package server

import (
	"github.com/Binaretech/classroom-auth/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func App() *fiber.App {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${method} ${path}\n",
	}))

	auth := app.Group("/auth")

	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/", handler.Verify)

	return app
}
