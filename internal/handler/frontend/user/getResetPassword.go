package user

import (
	"github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetResetPassword(c *echo.Context) error {

	config := config.Load()

	return response.Render(c, "frontend:user:/reset-password.html", map[string]any{
		"MinLengthPassword": config.Users.MinLengthPassword,
	})
}
