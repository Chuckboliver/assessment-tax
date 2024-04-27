package app

import (
	"log/slog"

	"github.com/chuckboliver/assessment-tax/common"
	"github.com/chuckboliver/assessment-tax/postgres"
	"github.com/chuckboliver/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func New() (*echo.Echo, error) {
	e := echo.New()
	db, err := postgres.New(postgres.PostgresConfig{})
	if err != nil {
		slog.Error("Failed to connect to postgres", err)
		return nil, err
	}

	taxConfigRepo := tax.NewTaxConfigPostgresRepository(db)
	taxCalculator := tax.NewCalculator(taxConfigRepo)
	taxController := tax.NewTaxController(taxCalculator)

	configureController(e, &taxController)

	return e, nil
}

func configureController(e *echo.Echo, controllers ...common.Controller) {
	for _, v := range controllers {
		v.RouteConfig(e)
	}
}
