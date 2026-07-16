package frontend

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (f *FrontendHandler) GetIndex(c *echo.Context) error {
	return response.Render(c, "frontend:index/index.html", map[string]any{})
}
