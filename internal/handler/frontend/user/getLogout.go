package users

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/model"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetLogout(c *echo.Context) error {
	auth := appAuth.Load[model.User, model.Session](c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "auth context not found")
	}

	_, err := auth.ActionLogout(c)
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/user/login")
}
