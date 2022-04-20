package server

import (
	"github.com/Binaretech/classroom-auth/errors"
	"github.com/Binaretech/classroom-auth/handler"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sirupsen/logrus"
)

// App holds the server and the routes to be used by the server to handle requests from the client side of the application
func App() *echo.Echo {
	app := echo.New()

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		println(err.Error())
		if err := errors.Handler(err, c); err != nil {
			logrus.Error(err)
		}
	}

	goth.UseProviders(
		google.New(
			"812826612012-pcaj0olcgshbbi4gfjhmj79svs5jkn5m.apps.googleusercontent.com",
			"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",
			"http://localhost/auth/google/callback",
			"email",
			"profile",
		),
	)

	auth := app.Group("/auth")

	auth.GET("/:provider/callback", func(c echo.Context) error {
		user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
		if err != nil {
			return err
		}

		return c.JSON(200, user)
	})

	auth.GET("/:provider", func(c echo.Context) error {
		gothic.BeginAuthHandler(c.Response(), c.Request())
		return nil
	})

	auth.POST("/login", handler.Login)
	auth.POST("/register", handler.Register)
	auth.POST("/refresh", handler.RefreshToken)
	auth.POST("/logout", handler.Logout)
	auth.POST("/", handler.Verify)

	return app
}
