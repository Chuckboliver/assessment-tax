package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/chuckboliver/assessment-tax/app"
	"github.com/chuckboliver/assessment-tax/common"
	"github.com/labstack/echo/v4"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable"
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "adminTax"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin!"
	}

	appConfig := common.AppConfig{
		Port:          port,
		DatabaseURL:   databaseURL,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
	}

	e, err := app.New(appConfig)
	if err != nil {
		slog.Error("Failed to create new echo server", err)
		os.Exit(1)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	address := fmt.Sprintf("0.0.0.0:%s", appConfig.Port)

	e.Logger.Fatal(e.Start(address))
}
