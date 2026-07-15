package user

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetLogin(c *echo.Context) error {
	return c.Render(http.StatusOK, "frontend:user:/login.html", map[string]any{})
}
