package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KKGo-Software-engineering/assessment-tax/common"
	"github.com/KKGo-Software-engineering/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	calc := tax.NewCalculator()
	result := calc.Calculate(tax.CalculationRequest{
		TotalIncome: 500000,
	})

	jsonStr, _ := json.Marshal(result)

	fmt.Printf("result: %v\n", string(jsonStr))

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	taxController := tax.NewTaxController(&calc)

	configureController(e, &taxController)

	e.Logger.Fatal(e.Start(":1323"))
}

func configureController(e *echo.Echo, controllers ...common.Controller) {
	for _, v := range controllers {
		v.RouteConfig(e)
	}
}
