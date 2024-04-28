package tax

import (
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
	group := e.Group("/tax/calculations")
	{
		group.POST("", c.calculateTax)
		group.POST("/upload-csv", c.calculateTaxFromUploadedCSV)
	}
}

type calculationRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

func (c *TaxController) calculateTax(ctx echo.Context) error {
	var request calculationRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Message: err.Error(),
		})
		return err
	}

	if err := ctx.Validate(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Message: err.Error(),
		})
		return err
	}

	result := c.taxCalculator.Calculate(ctx.Request().Context(), request)
	return ctx.JSON(http.StatusOK, result)
}

func (c *TaxController) calculateTaxFromUploadedCSV(ctx echo.Context) error {
	fileHeader, err := ctx.FormFile("taxFile")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Message: err.Error(),
		})
		return err
	}

	multipartFile, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Message: err.Error(),
		})
		return err
	}

	parser := newCSVParser()
	calculationRequests, err := parser.parseCalculationRequest(multipartFile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Message: err.Error(),
		})
		return err
	}

	result := c.taxCalculator.BatchCalculate(ctx.Request().Context(), calculationRequests)
	return ctx.JSON(http.StatusOK, result)
}
