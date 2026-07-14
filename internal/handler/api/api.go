package api

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/handler/api/users"
	"github.com/labstack/echo/v5"
)

func New(e *echo.Echo) {
	h := &ApiHandler{}

	api := e.Group("/api")
	api.POST("/lang", h.PostChangeLang)

	api.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})

	users.New(api)
}
