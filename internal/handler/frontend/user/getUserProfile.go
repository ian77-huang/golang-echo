package users

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetUserProfile(c *echo.Context) error {
	return c.Render(http.StatusOK, "frontend:users:/userProfile.html", map[string]any{})
}
