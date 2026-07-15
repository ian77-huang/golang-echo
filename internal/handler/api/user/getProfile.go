package user

import (
	"log"
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) GetProfile(c *echo.Context) error {
	log.Printf("\n==== 012345 ======\n")
	// h.DB
	user := auth.GetUser[model.User](c)

	userService := service.NewUserService(h.DB)

	userProfile, err := userService.GetUserProfile(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data":    userProfile,
		"message": shared.T(c, "user.profile.read.success"),
	})
}
