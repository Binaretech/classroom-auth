package server

import (
	"fmt"

	"github.com/Binaretech/classroom-auth/errors"
	"github.com/Binaretech/classroom-auth/handler"
	"github.com/Binaretech/classroom-auth/validation"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// App holds the server and the routes to be used by the server to handle requests from the client side of the application
func App(db *mongo.Database) *echo.Echo {
	app := echo.New()

	app.Validator = validation.SetUpValidator(db)

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		println(err.Error())
		if err := errors.Handler(err, c); err != nil {
			logrus.Error(err)
		}
	}

	auth := app.Group("/auth")

	handler := handler.NewHandler(db)

	handler.Routes(auth)

	auth.Any("/*", func(c echo.Context) error {
		req := map[string]any{}
		c.Bind(&req)

		fmt.Println(req)
		fmt.Println(c.Request().URL.Path)

		return c.JSON(404, "Not found")
	})

	return app
}
