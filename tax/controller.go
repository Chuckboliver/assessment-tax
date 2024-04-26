package tax

import (
	"fmt"
	"net/http"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/labstack/echo/v4"
)

var _ common.Controller = (*TaxController)(nil)

type TaxController struct {
	taxCalculator Calculator
}

func NewTaxController(taxCalculator Calculator) TaxController {
	return TaxController{
		taxCalculator: taxCalculator,
	}
}

func (c *TaxController) RouteConfig(e *echo.Echo) {
	e.POST("/tax/calculations", c.CalculateTax)
}

func (c *TaxController) CalculateTax(ctx echo.Context) error {
	var request CalculationRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	result := c.taxCalculator.Calculate(request)
	fmt.Printf("result: %v\n", result)
	return ctx.JSON(http.StatusOK, result)
}
