package admin

import (
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

type personalDeductionRequest struct {
	Amount float64 `json:"amount"`
}

func (a *AdminController) updatePersonalDeduction(ctx echo.Context) error {
	var request personalDeductionRequest
	if err := ctx.Bind(&request); err != nil {
		slog.Error("Failed to bind request", err)
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	err := a.adminService.UpdatePersonalDeduction(ctx.Request().Context(), request.Amount)
	if err != nil {
		slog.Error("Failed to update personal deduction", err)
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
