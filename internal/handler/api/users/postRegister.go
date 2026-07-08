package users

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

type Request struct {
	Account         string `json:"account" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

func PostRegister(c *echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	config := config.Load()
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

	// c.SetCookie(&http.Cookie{
	// 	Name:     "lang",
	// 	Value:    req.Account,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(365 * 24 * time.Hour),
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode,
	// })

	return response.JSON(c, map[string]any{
		"lang": req.Account,
	})
}
