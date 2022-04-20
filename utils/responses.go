package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ResponseBadRequest returns a 400 response with a message
func ResponseBadRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": msg,
	})
}

func ResponseError(c echo.Context, status int, err string) error {
	return c.JSON(status, map[string]string{
		"error": err,
	})
}
