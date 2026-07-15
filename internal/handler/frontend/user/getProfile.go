package user

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetProfile(c *echo.Context) error {
	return c.Render(http.StatusOK, "frontend:user:/profile.html", map[string]any{})
}
