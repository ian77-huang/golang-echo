package admin

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (f *adminHandler) GetIndex(c *echo.Context) error {
	return response.Render(c, "admin:index/index.html", map[string]any{})
}
