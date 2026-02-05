package middleware

import (
	"fmt"
	"net/http"
    customErrors "github.com/arunima10a/go-food-delivery/internal/common/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func CustomHTTPErrorHandler(logger zerolog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		if he, ok := err.(*customErrors.ApiError); ok {
			code = he.Status
			message = he.Message
		} else if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprint("%v", he.Message)
		}
		logger.Error().
			Err(err).
			Int("status", code).
			Str("method", c.Request().Method).
			Str("path", c.Path()).
			Msg("API Error Caught by Middleware")

		if !c.Response().Committed {
			c.JSON(code, customErrors.ApiError{
				Status:  code,
				Message: message,
			})
		}
	}
}
