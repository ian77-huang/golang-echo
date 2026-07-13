package users

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetIndex(c *echo.Context) error {
	return c.Render(http.StatusOK, "frontend:users:/index.html", map[string]any{})
}
