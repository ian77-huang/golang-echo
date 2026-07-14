package users

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetLogin(c *echo.Context) error {
	return c.Render(http.StatusOK, "frontend:users:/login.html", map[string]any{})
}
