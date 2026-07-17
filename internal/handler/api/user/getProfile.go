package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) GetProfile(c *echo.Context) error {
	user := auth.GetUser[model.User](c)
	if user == nil {
		return response.ErrorUnauthorized(c, "unauthorized")
	}
	userService := service.NewUserService(h.DB)

	userProfile, err := userService.GetUserProfile(user.ID)
	if err != nil {
		return err
	}
	if userProfile == nil {
		userProfile = &model.UserProfile{}
	}

	return response.JsonOk(c, map[string]any{
		"data":    userProfile,
		"message": shared.T(c, "user.profile.read.success"),
	})
}
