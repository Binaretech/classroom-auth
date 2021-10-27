package server

import (
	"github.com/Binaretech/classroom-auth/internal/handler"
	"github.com/gofiber/fiber/v2"
)

// App holds the server and the routes to be used by the server to handle requests from the client side of the application
func App() *fiber.App {
	app := fiber.New()

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/", handler.Verify)

	return app
}
