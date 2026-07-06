package users

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetRegister(c *echo.Context) error {
	return c.Render(http.StatusOK, "frontend:users:/register.html", map[string]any{
		"Name": "Yien12345",
	})
}
