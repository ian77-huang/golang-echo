package routing

import (
	"github.com/ian77-huang/golang-echo/internal/handler"
	handlerUsers "github.com/ian77-huang/golang-echo/internal/handler/frontend/users"
	"github.com/labstack/echo/v5"
)

func FrontendRouting(e *echo.Echo) {
	e.GET("/", handler.GetIndex)

	usersGroup := e.Group("/users")

	usersGroup.GET("", handlerUsers.GetIndex)
	usersGroup.GET("/login", handlerUsers.GetLogin)
	usersGroup.GET("/register", handlerUsers.GetRegister)
	usersGroup.GET("/logout", handlerUsers.GetLogout)
}
