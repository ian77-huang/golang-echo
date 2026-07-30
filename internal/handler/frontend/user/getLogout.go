package user

import (
	"log"

	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/model"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

func (h *UserHandler) GetLogout(c *echo.Context) error {
	auth := appAuth.Load[model.User, model.Session](c)
	if auth == nil {
		return response.ErrorInternalServerError(c, "auth context not found")
	}
	log.Printf("\n===== GetLogout 1 =====\n")
	_, err := auth.ActionLogout(c)
	log.Printf("\n===== GetLogout 2 =====\n")
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}
	log.Printf("\n===== GetLogout 3 =====\n")
	return response.Redirect(c, "/user/login")
}
