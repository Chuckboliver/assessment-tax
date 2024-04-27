package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/chuckboliver/assessment-tax/app"
	"github.com/labstack/echo/v4"
)

func main() {
	e, err := app.New()
	if err != nil {
		slog.Error("Failed to create new echo server", err)
		os.Exit(1)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
