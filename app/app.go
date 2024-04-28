package app

import (
	"log/slog"

	"github.com/chuckboliver/assessment-tax/admin"
	"github.com/chuckboliver/assessment-tax/common"
	"github.com/chuckboliver/assessment-tax/postgres"
	"github.com/chuckboliver/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func New(config common.AppConfig) (*echo.Echo, error) {
	db, err := postgres.New(config.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to postgres", err)
		return nil, err
	}

	taxConfigRepo := tax.NewTaxConfigPostgresRepository(db)
	taxCalculator := tax.NewCalculator(taxConfigRepo)
	taxController := tax.NewTaxController(taxCalculator)

	adminRepo := admin.NewAdminRepository(db)
	adminService := admin.NewAdminService(adminRepo)
	adminController := admin.NewAdminController(adminService, config)

	e := common.NewConfiguredEcho()

	configureController(e, &taxController, &adminController)

	return e, nil
}

func configureController(e *echo.Echo, controllers ...common.Controller) {
	for _, v := range controllers {
		v.RouteConfig(e)
	}
}
