package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetProfile(c *echo.Context) error {
	return response.Render(c, "frontend:user:/profile.html", map[string]any{})
}
