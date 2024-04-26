package main

import (
	"net/http"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/chuckboliver/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	calculator := tax.NewCalculator()
	taxController := tax.NewTaxController(&calculator)

	configureController(e, &taxController)

	e.Logger.Fatal(e.Start(":1323"))
}

func configureController(e *echo.Echo, controllers ...common.Controller) {
	for _, v := range controllers {
		v.RouteConfig(e)
	}
}
