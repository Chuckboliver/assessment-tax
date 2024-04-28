package tax

import (
	"log/slog"
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
	e.POST("/tax/calculations", c.calculateTax)
}

type calculationRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

func (c *TaxController) calculateTax(ctx echo.Context) error {
	var request calculationRequest
	if err := ctx.Bind(&request); err != nil {
		slog.Error("Failed to bind request", err)
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	if err := ctx.Validate(&request); err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	result := c.taxCalculator.Calculate(ctx.Request().Context(), request)
	return ctx.JSON(http.StatusOK, result)
}
