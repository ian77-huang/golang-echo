package users

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/config"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetRegister(c *echo.Context) error {

	config := config.Load()

	return c.Render(http.StatusOK, "frontend:users:/register.html", map[string]any{
		"MinLengthAccount":  config.Users.MinLengthAccount,
		"MinLengthPassword": config.Users.MinLengthPassword,
	})
}
