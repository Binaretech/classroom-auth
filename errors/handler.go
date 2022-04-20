package errors

import (
	"encoding/json"
	"net/http"

	"github.com/Binaretech/classroom-auth/lang"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Handler catch all the errors and returns a propper response based on the error type
func Handler(err error, c echo.Context) error {

	if e, ok := err.(*json.UnmarshalTypeError); ok {
		return response(c, NewBadRequestWrapped(lang.Trans("invalid data type"), e))
	}

	if e, ok := err.(ServerError); ok {
		return response(c, e)
	}

	return response(c, WrapError(err, lang.Trans("internal error"), http.StatusInternalServerError))
}

func response(c echo.Context, err ServerError) error {
	message := echo.Map{
		"error": err.GetMessage(),
	}

	if viper.GetBool("debug") {
		message["debug"] = err.Error()
	}

	return c.JSON(int(err.GetCode()), message)
}
