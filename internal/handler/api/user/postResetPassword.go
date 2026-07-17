package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) PostResetPassword(c *echo.Context) error {
	var req RequestResetPassword

	if err := c.Bind(&req); err != nil {
		return response.ErrorBadRequest(c, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	user := appAuth.GetUser[model.User](c)
	if user == nil {
		return response.ErrorInternalServerError(c, "auth context not found")
	}
	auth := appAuth.Load[model.User, model.Session](c)
	if auth == nil {
		return response.ErrorInternalServerError(c, "auth context not found")
	}

	updateUser, err := auth.ActionResetPassword(c, user.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return response.JsonOk(c, map[string]any{
		"message": shared.T(c, "user.auth.resetPasswordSuccess"),
		"id":      updateUser.ID,
	})
}
