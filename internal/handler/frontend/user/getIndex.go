package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetIndex(c *echo.Context) error {
	authUser := appAuth.GetUser[model.User](c)
	userService := service.NewUserService(h.DB)
	user, err := userService.GetUser(authUser.ID)
	if err != nil {
		return err
	}
	profile, err := userService.GetUserProfile(authUser.ID)
	if err != nil {
		return err
	}
	return response.Render(c, "frontend:user:/index.html", map[string]any{
		"user":    user.Data,
		"profile": profile,
	})
}
