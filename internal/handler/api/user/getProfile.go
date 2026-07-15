package user

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) GetProfile(c *echo.Context) error {
	user := auth.GetUser[model.User](c)

	userService := service.NewUserService(h.DB)

	userProfile, err := userService.GetUserProfile(user.ID)
	if err != nil {
		return err
	}
	if userProfile == nil {
		userProfile = &model.UserProfile{}
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data":    userProfile,
		"message": shared.T(c, "user.profile.read.success"),
	})
}
