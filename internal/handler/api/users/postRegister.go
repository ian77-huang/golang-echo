package users

import (
	"net/http"

	// "time"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/models/session"
	"github.com/ian77-huang/golang-echo/internal/models/users"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"

	// "github.com/ian77-huang/golang-echo/pkg/argon2"
	// "github.com/ian77-huang/golang-echo/pkg/database"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

type Request struct {
	Account         string `json:"account" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}
type Users struct {
}
type Session struct {
}

func PostRegister(c *echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
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
	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	auth := appAuth.GetAuth[users.User, session.Session](c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "auth context not found")
	}
	ID, err := auth.Register(c, req.Account, req.Password)
	if err != nil {
		return response.ValidationErrorAuth(c, err.(appAuth.FieldError))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": shared.T(c, "users.auth.create.success"),
		"id":      ID,
	})
}
