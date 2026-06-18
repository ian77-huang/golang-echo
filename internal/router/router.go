package router

import (
	"net/http"

	"github.com/labstack/echo/v5/middleware"
	// "myapp/internal/handler"

	"github.com/labstack/echo/v5"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// root route
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})

	api := e.Group("/api/v1")
	api.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})

	return e
}
