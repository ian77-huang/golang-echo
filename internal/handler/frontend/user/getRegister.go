package user

import (
	"github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetRegister(c *echo.Context) error {

	config := config.Load()

	return response.Render(c, "frontend:user:/register.html", map[string]any{
		"MinLengthAccount":  config.Users.MinLengthAccount,
		"MinLengthPassword": config.Users.MinLengthPassword,
	})
}
