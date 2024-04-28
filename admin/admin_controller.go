package admin

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/labstack/echo/v4"
)

var _ common.Controller = (*AdminController)(nil)

type AdminController struct {
	adminService AdminService
}

func NewAdminController(adminService AdminService) AdminController {
	return AdminController{
		adminService: adminService,
	}
}

func (a *AdminController) RouteConfig(e *echo.Echo) {
	e.POST("/admin/deductions/personal", a.updatePersonalDeduction)
}

type updatePersonalDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

type updatePersonalDeductionResponse struct {
	PersonalDeduction common.Float64 `json:"personalDeduction"`
}

func (a *AdminController) updatePersonalDeduction(ctx echo.Context) error {
	var request updatePersonalDeductionRequest
	if err := ctx.Bind(&request); err != nil {
		slog.Error("Failed to bind request", err)
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	if err := ctx.Validate(&request); err != nil {
		fmt.Printf("err: %v\n", err)
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	updatedPersonalDeduction, err := a.adminService.UpdatePersonalDeduction(ctx.Request().Context(), request.Amount)
	if err != nil {
		slog.Error("Failed to update personal deduction", err)
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	response := updatePersonalDeductionResponse{
		PersonalDeduction: common.Float64(updatedPersonalDeduction),
	}

	return ctx.JSON(http.StatusOK, response)
}
