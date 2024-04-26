package common

import "github.com/labstack/echo/v4"

type Controller interface {
	RouteConfig(e *echo.Echo)
}
