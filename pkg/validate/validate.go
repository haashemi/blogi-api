package validate

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type EchoValidator struct {
	validator *validator.Validate
}

func NewEchoValidator() *EchoValidator {
	return &EchoValidator{
		validator: validator.New(),
	}
}

func (cv *EchoValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
