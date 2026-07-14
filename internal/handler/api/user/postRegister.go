package users

import (
	"net/http"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"

	"github.com/ian77-huang/golang-echo/model"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) PostRegister(c *echo.Context) error {
	var req RequestRegister

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	config := appConfig.Load()
	if len(req.Account) < config.Users.MinLengthAccount {
		err := response.NewFieldError("account", "min", config.Users.MinLengthAccount)
		return response.ValidationCustomError(c, err)
	}
	if len(req.Password) < config.Users.MinLengthPassword {
		err := response.NewFieldError("password", "min", config.Users.MinLengthPassword)
		return response.ValidationCustomError(c, err)
	}

	auth := appAuth.Load[model.User, model.Session](c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "auth context not found")
	}

	ID, err := auth.ActionRegister(c, req.Account, req.Password)
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": shared.T(c, "users.auth.create.success"),
		"id":      ID,
	})
}
