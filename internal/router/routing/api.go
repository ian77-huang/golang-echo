package routing

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Api(e *echo.Echo) {
	api := e.Group("/api/v1")
	api.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})
}