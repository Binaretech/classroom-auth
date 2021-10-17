package server

import (
	"github.com/Binaretech/classroom-auth/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func App() *fiber.App {
	app := fiber.New()

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/", handler.Verify)

	return app
}
