package main

import (
	"net/http"

	"github.com/chuckboliver/assessment-tax/server"
	"github.com/labstack/echo/v4"
)

func main() {
	e := server.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
