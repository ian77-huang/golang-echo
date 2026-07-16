package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) PostLogin(c *echo.Context) error {
	var req RequestLogin

	if err := c.Bind(&req); err != nil {
		return response.ErrorBadRequest(c, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	auth := appAuth.Load[model.User, model.Session](c)
	if auth == nil {
		return response.ErrorInternalServerError(c, "auth context not found")
	}
	ID, err := auth.ActionLogin(c, req.Account, req.Password)
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return response.JsonOk(c, map[string]any{
		"message": shared.T(c, "user.auth.login.success"),
		"id":      ID,
	})
}
