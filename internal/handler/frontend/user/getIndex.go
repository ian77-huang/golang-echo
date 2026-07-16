package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetIndex(c *echo.Context) error {
	return response.Render(c, "frontend:user:/index.html", map[string]any{})
}
