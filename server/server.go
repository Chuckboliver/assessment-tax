package server

import (
	"github.com/chuckboliver/assessment-tax/common"
	"github.com/chuckboliver/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	taxCalculator := tax.NewCalculator()
	taxController := tax.NewTaxController(&taxCalculator)

	configureController(e, &taxController)

	return e
}

func configureController(e *echo.Echo, controllers ...common.Controller) {
	for _, v := range controllers {
		v.RouteConfig(e)
	}
}
