package server

import (
	"github.com/Binaretech/classroom-auth/errors"
	"github.com/Binaretech/classroom-auth/handler"
	"github.com/gofiber/fiber/v2"
)

// App holds the server and the routes to be used by the server to handle requests from the client side of the application
func App() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.Handler,
	})

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/refresh", handler.RefreshToken)
	auth.Post("/logout", handler.Logout)
	auth.Get("/", handler.Verify)

	return app
}
