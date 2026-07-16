package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/service"

	"github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) PutProfile(c *echo.Context) error {
	var req RequestPutProfile

	if err := c.Bind(&req); err != nil {
		return response.ErrorBadRequest(c, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	user := auth.GetUser[model.User](c)

	userService := service.NewUserService(h.DB)

	profile, err := userService.UpdateUserProfile(user.ID, &model.UserProfile{
		Name: req.Name, Email: req.Email, Phone: req.Phone, Bio: req.Bio,
	})
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return response.JsonOk(c, map[string]any{
		"message": shared.T(c, "user.profile.update.success"),
		"data":    profile,
	})
}
