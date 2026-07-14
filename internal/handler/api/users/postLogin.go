package users

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/models/session"
	"github.com/ian77-huang/golang-echo/internal/models/users"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) PostLogin(c *echo.Context) error {
	var req RequestLogin

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	auth := appAuth.Load[users.User, session.Session](c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "auth context not found")
	}

	ID, err := auth.ActionLogin(c, req.Account, req.Password)
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": shared.T(c, "users.auth.login.success"),
		"id":      ID,
	})
}
