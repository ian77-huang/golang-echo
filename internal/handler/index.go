package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetIndex(c *echo.Context) error {

	return c.Render(http.StatusOK, "frontend:index/index.html", map[string]string{"message": "Hello, World12!"})
}
