package common

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller interface {
	RouteConfig(e *echo.Echo)
}

func NewConfiguredEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &EchoValidator{Validator: validator.New()}
	return e
}
