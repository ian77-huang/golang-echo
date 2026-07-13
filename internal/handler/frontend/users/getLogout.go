package users

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/models/session"
	"github.com/ian77-huang/golang-echo/internal/models/users"
	"github.com/ian77-huang/golang-echo/internal/response"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func GetLogout(c *echo.Context) error {
	auth := appAuth.Load[users.User, session.Session](c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "auth context not found")
	}

	_, err := auth.ActionLogout(c)
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/users/login")
}
