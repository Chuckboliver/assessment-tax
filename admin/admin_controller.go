package admin

import (
	"log/slog"
	"net/http"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ common.Controller = (*AdminController)(nil)

type AdminController struct {
	adminService AdminService
	appConfig    common.AppConfig
}

func NewAdminController(adminService AdminService, appConfig common.AppConfig) AdminController {
	return AdminController{
		adminService: adminService,
		appConfig:    appConfig,
	}
}

func (a *AdminController) RouteConfig(e *echo.Echo) {
	group := e.Group("/admin/deductions")
	group.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == a.appConfig.AdminUsername && password == a.appConfig.AdminPassword, nil
	}))
	{
		group.POST("/personal", a.updatePersonalDeduction)
		group.POST("/k-receipt", a.updateKReceiptDeduction)
	}
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
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	if err := ctx.Validate(&request); err != nil {
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

type updateKReceiptDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

type updateKReceiptDeductionResponse struct {
	KReceipt common.Float64 `json:"kReceipt"`
}

func (a *AdminController) updateKReceiptDeduction(ctx echo.Context) error {
	var request updateKReceiptDeductionRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	if err := ctx.Validate(&request); err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}

	updatedKReceiptDeduction, err := a.adminService.UpdateKReceiptDeduction(ctx.Request().Context(), request.Amount)
	if err != nil {
		slog.Error("Failed to update personal deduction", err)
		ctx.NoContent(http.StatusInternalServerError)
		return err
	}

	response := updateKReceiptDeductionResponse{
		KReceipt: common.Float64(updatedKReceiptDeduction),
	}

	return ctx.JSON(http.StatusOK, response)
}
