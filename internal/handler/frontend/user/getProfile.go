package user

import (
	appConfig "github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetProfile(c *echo.Context) error {
	config := appConfig.Load()
	maxSize := int64(config.MaxSizeUserProfileAvatar)

	return response.Render(c, "frontend:user:/profile.html", map[string]any{
		"MaxSize": maxSize,
	})
}
