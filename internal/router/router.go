package router

import (
	"github.com/ian77-huang/golang-echo/internal/handler"
	"github.com/ian77-huang/golang-echo/internal/router/routing"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// root route
	e.GET("/", handler.GetIndex)

	routing.Api(e)

	return e
}
