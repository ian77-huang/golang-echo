package routing

import (
	"net/http"

	apiHandler "github.com/ian77-huang/golang-echo/internal/handler/api"
	"github.com/labstack/echo/v5"
)

func ApiRouting(e *echo.Echo) {
	api := e.Group("/api")
	api.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})
	api.POST("/lang", apiHandler.PostChangeLang)
}
